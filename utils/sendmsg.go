package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func SendText(chatId int64, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, text)
	bot.Send(msg)
}

func SendMessage(bot *tgbotapi.BotAPI, c tgbotapi.Chattable, fmt string, args ...interface{}) {
	if _, err := bot.Send(c); err != nil {
		log.Printf(fmt, args...)
		log.Println(err)
	}
}
