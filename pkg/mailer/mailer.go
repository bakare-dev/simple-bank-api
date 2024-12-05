package mailer

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"

	"github.com/bakare-dev/simple-bank-api/pkg/config"
)

type Mailer struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
}

func NewMailer() *Mailer {
	return &Mailer{
		smtpHost:     config.Settings.Infrastructure.Mailer.SmtpHost,
		smtpPort:     config.Settings.Infrastructure.Mailer.SmtpPort,
		smtpUser:     config.Settings.Infrastructure.Mailer.SmtpUser,
		smtpPassword: config.Settings.Infrastructure.Mailer.SmtpPassword,
	}
}

type MailInfo struct {
	Sender       string
	Recipients   []string
	Subject      string
	TemplateFile string
	Data         interface{}
}

func (m *Mailer) SendMail(info MailInfo) error {
	tmpl, err := template.ParseFiles(info.TemplateFile)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, info.Data); err != nil {
		log.Printf("Error executing template: %v", err)
		return err
	}

	message := "From: " + info.Sender + "\r\n" +
		"To: " + joinRecipients(info.Recipients) + "\r\n" +
		"Subject: " + info.Subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body.String()

	auth := smtp.PlainAuth("", m.smtpUser, m.smtpPassword, m.smtpHost)

	err = smtp.SendMail(m.smtpHost+":"+m.smtpPort, auth, info.Sender, info.Recipients, []byte(message))
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Println("Email sent successfully to", info.Recipients)
	return nil
}

func joinRecipients(recipients []string) string {
	var buffer bytes.Buffer
	for i, recipient := range recipients {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(recipient)
	}
	return buffer.String()
}
