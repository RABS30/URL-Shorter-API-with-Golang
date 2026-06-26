package helper

import (
	"context"
	"net/smtp"
	"shorter-url/internal/domain"
)

type smtpEmailSender struct {
	host string
	port string
	auth smtp.Auth
	from string
}

func NewEmailSender(host, port, email, password string) domain.EmailSender {
	auth := smtp.PlainAuth("", email, password, host)

	return &smtpEmailSender{
		host: host,
		port: port,
		auth: auth,
		from: email,
	}
}

func (s *smtpEmailSender) SendEmail(ctx context.Context, to, subject, body string) error {
	message := []byte(
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body + "\r\n")

	addrs := s.host + ":" + s.port
	return smtp.SendMail(addrs, s.auth, s.from, []string{to}, message)
}
