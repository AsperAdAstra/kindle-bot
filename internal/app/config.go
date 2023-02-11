package app

import (
	"github.com/AsperAdAstra/kindle-bot/internal/mailer"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token    string `required:"true" envconfig:"BOT_TOKEN"`
	User     int64  `required:"true" envconfig:"USER_ID"`
	MailConf mailer.MailConfig
	SMTPConf mailer.SmtpConfig
}

func NewConfig() *Config {
	c := &Config{}
	err := envconfig.Process("", c)
	if err != nil {
		panic(err)
	}
	return c
}
