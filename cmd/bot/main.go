package main

import (
	"github.com/AsperAdAstra/kindle-bot/internal/app"
	"github.com/AsperAdAstra/kindle-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	Config := app.NewConfig()

	Bot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		panic(err)
	}

	Bot.Debug = true
	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	Handler := handler.NewHandler(Config, Bot)

	updates := Bot.GetUpdatesChan(u)
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		log.Printf("Data folder does not exist. Creating...")
		if err := os.Mkdir("./data", 0755); err != nil {
			log.Fatal(err)
		}
	}

	for update := range updates {
		if err := Handler.HandleIncomingMessage(update); err != nil {
			log.Printf("Error: %s", err.Error())
		}
	}
}
