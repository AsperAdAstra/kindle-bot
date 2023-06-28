package handler

import (
	"errors"
	"fmt"
	"github.com/AsperAdAstra/kindle-bot/internal/app"
	botHelper "github.com/AsperAdAstra/kindle-bot/internal/bot"
	"github.com/AsperAdAstra/kindle-bot/internal/mailer"
	"github.com/AsperAdAstra/kindle-bot/internal/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bolt "go.etcd.io/bbolt"
	"log"
	"strings"
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
	db     *bolt.DB
}

func NewHandler(config *app.Config, bot *tgbotapi.BotAPI, db *bolt.DB) *Handler {
	return &Handler{
		config: config,
		bot:    bot,
		dest:   config.DataDir,
		db:     db,
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

	if update.Message.IsCommand() {
		return h.handleCommand(update)
	}

	doc := update.Message.Document
	if doc == nil {
		return errNoFile
	}

	if !h.bookExists(doc.FileName) {
		dest, err := h.downloadFile(doc)
		if err != nil {
			return errors.New("download failed")
		}
		log.Printf("File downloaded to: %s", dest)

		if err := h.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("books"))
			if b == nil {
				b, err = tx.CreateBucket([]byte("books"))
				if err != nil {
					return err
				}
			}
			if err := b.Put([]byte(doc.FileName), []byte(dest)); err != nil {
				log.Printf("Error: %s", err.Error())
				return err
			}
			return nil
		}); err != nil {
			return err
		}

		return h.forwardBook(update, doc.FileName)
	}
	return errors.New("already forwarded")
}

func (h Handler) handleCommand(update tgbotapi.Update) error {
	var msg tgbotapi.MessageConfig
	switch update.Message.Command() {
	case "start":
		msg = botHelper.NewMessageReply("Hello!", update.Message.Chat.ID, update.Message.MessageID)
	case "help":
		msg = botHelper.NewMessageReply("Send me a book and I'll forward it to your Kindle. \r\n Send /list to see already forwarded books.", update.Message.Chat.ID, update.Message.MessageID)
	case "list":
		books, err := h.ListBooks()
		if err != nil {
			return err
		}
		msg = botHelper.NewMessageReply(fmt.Sprintf("List of forwarded books: \r\n%s", strings.Join(books, "\r\n")), update.Message.Chat.ID, update.Message.MessageID)
	}
	_, err := h.bot.Send(msg)
	return err
}

func (h Handler) ListBooks() ([]string, error) {
	var list []string

	err := h.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("books"))
		if b == nil {
			return nil
		}
		if err := b.ForEach(func(k, v []byte) error {
			list = append(list, string(k))
			return nil
		}); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (h Handler) forwardBook(update tgbotapi.Update, dest string) error {
	if h.config.DryRun {
		log.Printf("Dry run: email not sent to %s", h.config.MailConf.To)
		msg := botHelper.NewMessageReply("dry run", update.Message.Chat.ID, update.Message.MessageID)
		if _, err := h.bot.Send(msg); err != nil {
			return err
		}
		return nil
	}

	// Send email
	mailMsg := mailer.Compose(h.config.MailConf.From, h.config.MailConf.To, update.Message.Document.FileName, dest)
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

func (h Handler) bookExists(name string) bool {
	err := h.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("books"))
		if b == nil {
			return nil
		}

		if b.Get([]byte(name)) != nil {
			return errors.New("already forwarded")
		}

		return nil
	})

	if err != nil {
		return true
	}
	return false
}

func (h Handler) prepareFilename(filename string) (string, error) {
	if filename == "" { // If doc is not present in message
		return "", errNoAttachment
	}

	dest := fmt.Sprintf("%s/%s", h.dest, filename)
	return dest, nil
}

func (h Handler) downloadFile(doc *tgbotapi.Document) (string, error) {
	fileUrl, err := h.bot.GetFileDirectURL(doc.FileID)
	if err != nil {
		return "", err
	}

	dest, err := h.prepareFilename(doc.FileName)
	if err != nil {
		return "", err
	}

	if err := transport.DownloadFile(fileUrl, dest); err != nil {
		return "", err
	}

	return dest, nil
}
