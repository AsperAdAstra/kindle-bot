package mailer

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
