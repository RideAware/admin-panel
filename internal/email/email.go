package email

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/rideaware/admin-panel/internal/config"
	"github.com/rideaware/admin-panel/internal/database"

	"github.com/wneessen/go-mail"
)

// SendUpdate sends a newsletter with the given subject and body to all subscriber emails stored in the database.
// It returns a human-readable status message and, when subscriber retrieval fails, the underlying error.
// - If retrieving subscribers fails: returns "Failed to retrieve subscribers" and the error.
// - If no subscribers are found: returns "No subscribers found." and nil.
// - If sending to a specific subscriber fails: returns "Failed to send to <email>" and nil.
// - On success: returns "Email has been sent to all subscribers." and nil.
// Note: logging the newsletter entry in the database is attempted after sending and any logging failure is non-fatal.
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

// send constructs and sends an HTML newsletter update to the specified recipient using the current SMTP configuration.
// It embeds an unsubscribe link for the recipient and returns true if the message was sent successfully, false if client creation, message setup, or sending fails.
func send(subject, body, recipient string) bool {
	cfg := config.Current

	var opts []mail.ClientOption
	opts = append(opts,
		mail.WithPort(cfg.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.SMTPUser),
		mail.WithPassword(cfg.SMTPPassword),
		mail.WithTimeout(10*time.Second),
	)

	// Use SSL for port 465, STARTTLS for others
	if cfg.SMTPPort == 465 {
		opts = append(opts, mail.WithSSL())
	} else {
		opts = append(opts, mail.WithTLSPolicy(mail.TLSMandatory))
	}

	client, err := mail.NewClient(cfg.SMTPServer, opts...)
	if err != nil {
		log.Printf("Failed to create mail client: %v", err)
		return false
	}
	defer client.Close()

	m := mail.NewMsg()
	if err := m.From(cfg.SenderEmail); err != nil {
		log.Printf("Failed to set from: %v", err)
		return false
	}
	if err := m.To(recipient); err != nil {
		log.Printf("Failed to set to: %v", err)
		return false
	}
	m.Subject(subject)

	unsubLink := fmt.Sprintf("https://%s/unsubscribe?email=%s",
		cfg.BaseURL, url.QueryEscape(recipient))

	// Build HTML body with unsubscribe link
	htmlBody := buildHTMLBody(body, unsubLink)
	m.SetBodyString(mail.TypeTextHTML, htmlBody)

	if err := client.Send(m); err != nil {
		log.Printf("Failed to send email to %s: %v", recipient, err)
		return false
	}

	log.Printf("Update email sent to: %s", recipient)
	return true
}

// buildHTMLBody constructs the final HTML email body by appending an unsubscribe footer to the user-provided content.
// It handles both complete HTML documents and HTML fragments.
func buildHTMLBody(body, unsubLink string) string {
	footer := fmt.Sprintf(
		"<br><br><hr><p style='font-size: 12px; color: #666;'>If you ever wish to unsubscribe, "+
			"please click <a href='%s'>here</a>.</p>",
		unsubLink)

	// If body contains closing html tag, insert before it
	if strings.Contains(strings.ToLower(body), "</html>") {
		return strings.Replace(body, "</html>", footer+"</html>", 1)
	}

	// If body contains closing body tag, insert before it
	if strings.Contains(strings.ToLower(body), "</body>") {
		return strings.Replace(body, "</body>", footer+"</body>", 1)
	}

	// Otherwise just append
	return body + footer
}