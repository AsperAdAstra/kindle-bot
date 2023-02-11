package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewMessageReply(text string, chatId int64, messageId int) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyToMessageID = messageId
	return msg
}
