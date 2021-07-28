package emailservice

import (
	"net/smtp"
	"os"
)

func NewEmailService() *EmailService {
	from := os.Getenv("MAIL_FROM")
	pass := os.Getenv("MAIL_PASS")
	host := "smtp.gmail.com"
	port := "587"
	addr := host + ":" + port
	auth := smtp.PlainAuth("", from, pass, host)

	return &EmailService{
		Auth: auth,
		Port: port,
		Addr: addr,
		From: from,
	}
}

type EmailService struct {
	Auth smtp.Auth
	Port string
	Addr string
	From string
}

func (es EmailService) Send(to, subject, html string) error {
	sub := "Subject: " + subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(sub + mime + html)
	return smtp.SendMail(es.Addr, es.Auth, es.From, []string{to}, msg)
}
