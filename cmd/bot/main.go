package main

import (
	"github.com/AsperAdAstra/kindle-bot/internal/app"
	"github.com/AsperAdAstra/kindle-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	Config := app.NewConfig()

	Bot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		panic(err)
	}

	Bot.Debug = Config.Debug
	log.Printf("Authorized on account %s", Bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = Config.Timeout
	Handler := handler.NewHandler(Config, Bot)

	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		log.Printf("Data folder does not exist. Creating...")
		if err := os.Mkdir("./data", 0755); err != nil {
			log.Fatal(err)
		}
	}

	updates := Bot.GetUpdatesChan(u)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for update := range updates {
			if err := Handler.HandleIncomingMessage(update); err != nil {
				log.Printf("Error: %s", err.Error())
			}
		}
	}()
	log.Println("Bot started. Press Ctrl+C to stop")
	defer func() {
		log.Println("Bot stopped")
	}()
	<-c
}
