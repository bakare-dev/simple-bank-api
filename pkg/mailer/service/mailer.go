package service

import (
	"log"
	"sync"

	"github.com/bakare-dev/simple-bank-api/pkg/mailer"
)

type NotificationService struct {
	mailer mailer.Mailer
}

func NewNotificationService(mailer mailer.Mailer) *NotificationService {
	return &NotificationService{mailer: mailer}
}

func (ns *NotificationService) SendVerifyRegistration(message Message) []Response {
	var wg sync.WaitGroup
	responses := make(chan Response, len(message.Recipients))

	for _, recipient := range message.Recipients {
		wg.Add(1)
		go func(recipient string) {
			defer wg.Done()

			info := mailer.MailInfo{
				Sender:       "noreply@bakaredev.me",
				TemplateFile: "../template/verify-registration.html",
				Subject:      "Account Created Successfully",
				Recipients:   []string{recipient},
				Data:         message.Data,
			}

			err := ns.mailer.SendMail(info)
			if err != nil {
				log.Printf("Failed to send email to %s: %v", recipient, err)
				responses <- Response{Status: "failed", Message: err.Error()}
			} else {
				responses <- Response{Status: "success", Message: "Email sent successfully to " + recipient}
			}
		}(recipient)
	}

	wg.Wait()
	close(responses)

	var result []Response
	for resp := range responses {
		result = append(result, resp)
	}

	return result
}

type Message struct {
	Recipients []string
	Data       interface{}
}

type Response struct {
	Status  string
	Message string
}
