package handler

import (
	"errors"
	"fmt"
	"github.com/AsperAdAstra/kindle-bot/internal/app"
	botHelper "github.com/AsperAdAstra/kindle-bot/internal/bot"
	"github.com/AsperAdAstra/kindle-bot/internal/mailer"
	"github.com/AsperAdAstra/kindle-bot/internal/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var (
	errNoMessage    = errors.New("no message")
	errUnauthorized = errors.New("unauthorized")
	errNoAttachment = errors.New("no attachment")
	errNoFile       = errors.New("no file")
)

type Handler struct {
	config *app.Config
	dest   string
	bot    *tgbotapi.BotAPI
}

func NewHandler(config *app.Config, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		config: config,
		bot:    bot,
		dest:   "data", // Yet hardcoded
	}
}

// Handle any received message, filter out
func (h Handler) HandleIncomingMessage(update tgbotapi.Update) error {
	if update.Message == nil { // If there is no message
		return errNoMessage
	}
	if update.Message.From.ID != h.config.User { // Ignore anyone except blocked user
		return errUnauthorized
	}

	doc := update.Message.Document
	dest, err := h.downloadFile(doc)
	if err != nil {
		log.Printf("[%s] %s. Response: nothing to forward", update.Message.From.UserName, update.Message.Text)
		msg := botHelper.NewMessageReply("nothing to forward", update.Message.Chat.ID, update.Message.MessageID)
		if _, err := h.bot.Send(msg); err != nil {
			return err
		}
		return errNoFile
	}

	// Send email
	mailMsg := mailer.Compose(h.config.MailConf.From, h.config.MailConf.To, doc.FileName, dest)
	if err := mailer.Send(h.config.SMTPConf, mailMsg); err != nil {
		return err
	}

	log.Printf("Email sent to %s", h.config.MailConf.To)
	msg := botHelper.NewMessageReply("sent", update.Message.Chat.ID, update.Message.MessageID)
	if _, err := h.bot.Send(msg); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Handle update.Message.Document -> download file do h.dest or return error
func (h Handler) downloadFile(doc *tgbotapi.Document) (string, error) {
	if doc == nil { // If doc is not present in message
		return "", errNoAttachment
	}

	fileUrl, err := h.bot.GetFileDirectURL(doc.FileID)
	if err != nil {
		return "", err
	}

	// Output folder is hardcoded -> ./data
	dest := fmt.Sprintf("%s/%s", h.dest, doc.FileName)
	if err := transport.DownloadFile(fileUrl, dest); err != nil {
		return "", err
	}

	return dest, nil
}
