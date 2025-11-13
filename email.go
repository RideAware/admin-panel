package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wneessen/go-mail"
)

func sendUpdateEmail(subject, body, recipient string) bool {
	client, err := mail.NewClient(
		config.SMTPServer,
		mail.WithPort(config.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.SMTPUser),
		mail.WithPassword(config.SMTPPassword),
		mail.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Printf("Failed to create mail client: %v", err)
		return false
	}
	defer client.Close()

	m := mail.NewMsg()
	if err := m.From(config.SenderEmail); err != nil {
		log.Printf("Failed to set from: %v", err)
		return false
	}
	if err := m.To(recipient); err != nil {
		log.Printf("Failed to set to: %v", err)
		return false
	}
	m.Subject(subject)

	unsubLink := fmt.Sprintf(
		"https://%s/unsubscribe?email=%s",
		config.BaseURL,
		recipient,
	)

	htmlBody := fmt.Sprintf(
		"%s<br><br>If you ever wish to unsubscribe, "+
			"please click <a href='%s'>here</a>",
		body, unsubLink,
	)
	m.SetBodyString(mail.TypeTextHTML, htmlBody)

	if err := client.Send(m); err != nil {
		log.Printf("Failed to send email to %s: %v", recipient, err)
		return false
	}

	log.Printf("Update email sent to: %s", recipient)
	return true
}

func processSendUpdateEmail(subject, body string) (string, error) {
	subscribers, err := getAllEmails()
	if err != nil {
		return "Failed to retrieve subscribers", err
	}

	if len(subscribers) == 0 {
		return "No subscribers found.", nil
	}

	for _, email := range subscribers {
		if !sendUpdateEmail(subject, body, email) {
			return fmt.Sprintf("Failed to send to %s", email),
				nil
		}
	}

	// Log newsletter
	_, err = db.Exec(
		"INSERT INTO newsletters (subject, body) VALUES ($1, $2)",
		subject, body,
	)
	if err != nil {
		log.Printf("Error logging newsletter: %v", err)
	}

	return "Email has been sent to all subscribers.", nil
}