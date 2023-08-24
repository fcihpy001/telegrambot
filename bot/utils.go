package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// test

func (bot *SmartBot) CheckAdmin(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	replyTo := msg.ReplyToMessage
	if chat == nil || replyTo == nil {
		log.Println("chat is nil or replyTo is nil")
		return
	}
	bot.CheckUserIsAdmin(chat.ID, replyTo.From.ID)
}
