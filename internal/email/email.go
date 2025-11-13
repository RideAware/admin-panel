package email

import (
	"fmt"
	"log"
	"time"

	"github.com/rideaware/admin-panel/internal/config"
	"github.com/rideaware/admin-panel/internal/database"

	"github.com/wneessen/go-mail"
)

func SendUpdate(subject, body string) (string, error) {
	subscribers, err := database.GetAllEmails()
	if err != nil {
		return "Failed to retrieve subscribers", err
	}

	if len(subscribers) == 0 {
		return "No subscribers found.", nil
	}

	for _, email := range subscribers {
		if !send(subject, body, email) {
			return fmt.Sprintf("Failed to send to %s", email), nil
		}
	}

	if err := database.LogNewsletter(subject, body); err != nil {
		log.Printf("Error logging newsletter: %v", err)
	}

	return "Email has been sent to all subscribers.", nil
}

func send(subject, body, recipient string) bool {
	cfg := config.Current

	client, err := mail.NewClient(
		cfg.SMTPServer,
		mail.WithPort(cfg.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.SMTPUser),
		mail.WithPassword(cfg.SMTPPassword),
		mail.WithTimeout(10*time.Second),
	)
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
		cfg.BaseURL, recipient)

	htmlBody := fmt.Sprintf(
		"%s<br><br>If you ever wish to unsubscribe, "+
			"please click <a href='%s'>here</a>",
		body, unsubLink)
	m.SetBodyString(mail.TypeTextHTML, htmlBody)

	if err := client.Send(m); err != nil {
		log.Printf("Failed to send email to %s: %v", recipient, err)
		return false
	}

	log.Printf("Update email sent to: %s", recipient)
	return true
}