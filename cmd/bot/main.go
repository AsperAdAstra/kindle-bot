package main

import (
	"fmt"
	"github.com/AsperAdAstra/kindle-bot/internal/app"
	botHelper "github.com/AsperAdAstra/kindle-bot/internal/bot"
	"github.com/AsperAdAstra/kindle-bot/internal/mailer"
	"github.com/AsperAdAstra/kindle-bot/internal/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	Config := app.NewConfig()

	bot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // If there is no message
			continue
		}
		if update.Message.From.ID != Config.User { // Ignore anyone except blocked user
			continue
		}

		if update.Message.Document == nil { // If document is not present in message
			log.Printf("[%s] %s. Response: nothing to forward", update.Message.From.UserName, update.Message.Text)
			msg := botHelper.NewMessageReply("nothing to forward", update.Message.Chat.ID, update.Message.MessageID)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
			continue
		}

		doc := update.Message.Document
		fileUrl, err := bot.GetFileDirectURL(doc.FileID)
		if err != nil {
			log.Fatal(err)
		}

		// Output folder is hardcoded -> ./data
		dest := fmt.Sprintf("%s/%s", "data", doc.FileName)
		if err := transport.DownloadFile(fileUrl, dest); err != nil {
			log.Fatal(err)
		}

		// Send email
		mailMsg := mailer.Compose(Config.MailConf.From, Config.MailConf.To, doc.FileName, dest)
		if err := mailer.Send(Config.SMTPConf, mailMsg); err != nil {
			log.Fatal(err)
		}

		log.Printf("Email sent to %s", Config.MailConf.To)
		msg := botHelper.NewMessageReply("sent", update.Message.Chat.ID, update.Message.MessageID)
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}
