package notifications

import (
	"fmt"
	"net/smtp"
	"strconv"

	cfg "github.com/akhilckenshi/notification/pkg/settings"
)

// SendEmail sends an email notification using SMTP or a third-party service.
func SendEmail(to string, subject string, body string) error {
	from := cfg.Config.AppEmailID
	password := cfg.Config.AppEmailPassword
	smtpHost := cfg.Config.SMTPHost
	smtpPort := strconv.Itoa(cfg.Config.SMTPPort)

	// Setup the authentication information.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Create the email message.
	// Notice the "Content-Type" header for HTML.
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=UTF-8\n\n%s", from, to, subject, body)

	// Send the email.
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}

// // SendEmail sends an email notification using SMTP or a third-party service.
// func SendEmail(to string, subject string, body string) error {
// 	username := cfg.Config.Email.Username
// 	from := cfg.Config.Email.Id
// 	password := cfg.Config.Email.Pwd
// 	smtpHost := cfg.Config.Email.SmtpHost
// 	smtpPort := strconv.Itoa(cfg.Config.Email.SmtpPort)
// 	// Setup the authentication information.
// 	auth := smtp.PlainAuth("", username, password, smtpHost)
// 	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
// 	if err != nil {
// 		fmt.Printf("Failed to send email to %s: %v\n", to, err)
// 		return err
// 	}

// 	fmt.Println("Email sent successfully!")
// 	return nil
// }
