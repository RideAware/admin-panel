package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/rideaware/admin-panel/internal/config"
	"github.com/rideaware/admin-panel/internal/database"
)

// SendUpdate sends a newsletter with the given subject and body to all subscriber emails stored in the database.
// It returns a human-readable status message and, when subscriber retrieval fails, the underlying error.
func SendUpdate(subject, body string) (string, error) {
	subscribers, err := database.GetAllEmails()
	if err != nil {
		return "Failed to retrieve subscribers", err
	}
	if len(subscribers) == 0 {
		return "No subscribers found.", nil
	}

	var succeeded, failed int
	for _, email := range subscribers {
		if send(subject, body, email) {
			succeeded++
		} else {
			failed++
		}
	}

	if err := database.LogNewsletter(subject, body); err != nil {
		log.Printf("Error logging newsletter: %v", err)
	}

	if failed == 0 {
		return fmt.Sprintf("Email sent to all %d subscribers.", succeeded), nil
	}
	return fmt.Sprintf("Sent to %d/%d subscribers; %d failed.", succeeded, succeeded+failed, failed), nil
}

// send constructs and sends an HTML newsletter update to the specified recipient
func send(subject, body, recipient string) bool {
	cfg := config.Current

	log.Printf("Attempting to send email to %s via %s:%d", recipient, cfg.SMTPServer, cfg.SMTPPort)

	addr := fmt.Sprintf("%s:%d", cfg.SMTPServer, cfg.SMTPPort)

	unsubLink := fmt.Sprintf("https://%s/unsubscribe?email=%s",
		cfg.BaseURL, url.QueryEscape(recipient))

	htmlBody := buildHTMLBody(body, unsubLink)
	message := buildMessage(cfg.SenderEmail, recipient, subject, htmlBody)

	// Create TLS connection
	tlsconfig := &tls.Config{
		ServerName: cfg.SMTPServer,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		log.Printf("Failed to connect to SMTP %s: %v", addr, err)
		return false
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, cfg.SMTPServer)
	if err != nil {
		log.Printf("Failed to create SMTP client: %v", err)
		return false
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPServer)
	if err := client.Auth(auth); err != nil {
		log.Printf("SMTP auth failed for %s: %v", cfg.SMTPUser, err)
		return false
	}

	// Send the email
	if err := client.Mail(cfg.SenderEmail); err != nil {
		log.Printf("MAIL command failed: %v", err)
		return false
	}

	if err := client.Rcpt(recipient); err != nil {
		log.Printf("RCPT command failed for %s: %v", recipient, err)
		return false
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("DATA command failed: %v", err)
		return false
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write message: %v", err)
		return false
	}

	err = w.Close()
	if err != nil {
		log.Printf("Failed to close DATA: %v", err)
		return false
	}

	client.Quit()

	log.Printf("Update email sent to: %s", recipient)
	return true
}

func buildMessage(from, to, subject, body string) string {
	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	msg += "\r\n"
	msg += body
	return msg
}

// buildHTMLBody constructs the final HTML email body by appending an unsubscribe footer to the user-provided content.
func buildHTMLBody(body, unsubLink string) string {
	footer := fmt.Sprintf(
		"<br><br><hr><p style='font-size: 12px; color: #666;'>If you ever wish to unsubscribe, "+
			"please click <a href='%s'>here</a>.</p>",
		unsubLink)

	if strings.Contains(strings.ToLower(body), "</html>") {
		return strings.Replace(body, "</html>", footer+"</html>", 1)
	}

	if strings.Contains(strings.ToLower(body), "</body>") {
		return strings.Replace(body, "</body>", footer+"</body>", 1)
	}

	return body + footer
}
