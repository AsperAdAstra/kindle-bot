package mailer

import (
	"github.com/go-mail/mail"
)

type SmtpConfig struct {
	Host     string `envconfig:"SMTP_HOST" required:"true"`
	Port     int    `envconfig:"SMTP_PORT" default:"587"`
	Username string `envconfig:"SMTP_USERNAME" required:"true""`
	Password string `envconfig:"SMTP_PASSWORD" required:"true"`
}

type MailConfig struct {
	From string `envconfig:"MAIL_FROM" required:"true"`
	To   string `envconfig:"MAIL_TO" required:"true"`
}

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
