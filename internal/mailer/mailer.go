package mailer

import (
	"github.com/go-mail/mail"
)

func Compose(from, to, subject, attachment string) *mail.Message {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", "This is the plain text.")
	m.Attach(attachment)

	return m
}

func Send(c SmtpConfig, m *mail.Message) error {
	d := mail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	return d.DialAndSend(m)
}
