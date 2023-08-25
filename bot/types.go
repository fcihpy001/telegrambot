package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type SmartBot struct {
	Token string
	Debug bool
	bot   *tgbotapi.BotAPI
}
