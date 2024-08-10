package email

import (
	"net/smtp"
	"os"
)

type EmailService struct {
    AuthInfo smtp.Auth
    Sender string
} 

type IEmailService interface {
    SendEmail(to []string, message []byte) error
    Authenticate() 
}

func(e *EmailService) Authenticate() {
    email := os.Getenv("SENDER_EMAIL")
    pass := os.Getenv("SENDER_EMAIL_PASS")
    e.Sender = email
    e.AuthInfo = smtp.PlainAuth("", email, pass, "smtp.gmail.com")
}

func(e *EmailService) SendEmail(to []string, message []byte) error {
  return smtp.SendMail("smtp.gmail.com:587", e.AuthInfo, e.Sender, to, message)
}
