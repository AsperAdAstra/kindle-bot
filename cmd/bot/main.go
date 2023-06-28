package main

import (
	"fmt"
	"github.com/AsperAdAstra/kindle-bot/internal/app"
	botHelper "github.com/AsperAdAstra/kindle-bot/internal/bot"
	"github.com/AsperAdAstra/kindle-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	Config := app.NewConfig()

	Bot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		panic(err)
	}

	Db, err := bolt.Open(fmt.Sprintf("%s/kindle-bot.db", Config.DataDir), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	Bot.Debug = Config.Debug
	log.Printf("Authorized on account %s", Bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = Config.Timeout
	Handler := handler.NewHandler(Config, Bot, Db)

	if _, err := os.Stat(Config.DataDir); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist. Creating...", Config.DataDir)
		if err := os.Mkdir(Config.DataDir, 0755); err != nil {
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
				msg := botHelper.NewMessageReply(err.Error(), update.Message.Chat.ID, update.Message.MessageID)
				if _, err := Bot.Send(msg); err != nil {
					log.Printf("Error: %s", err.Error())
				}
			}
		}
	}()
	log.Println("Bot started. Press Ctrl+C to stop")
	defer func() {
		log.Println("Bot stopped")
	}()
	<-c
}
