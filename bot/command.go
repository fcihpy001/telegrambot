package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/setting"
)

// 处理以/开头的指令消息,如/help  /status等
func (bot *SmartBot) handleCommand(update tgbotapi.Update) {
	fmt.Println("command---", update.Message.Command())
	switch strings.ToLower(update.Message.Command()) {
	case "help":
		setting.Help(update.Message.Chat.ID, bot.bot)
	case "start":
		setting.Settings(update.Message.Chat.ID, bot.bot)

	case "create":

	case "luck":

	case "stat":

	case "stats":

	case "stat_week":

	case "filter":

	case "stop":

	case "filters":

	case "mute":

	case "unmute":

	case "ban":
		//bot.Ban(update)
	case "unban":
		//bot.UnBan(update)
	case "admin":
		//bot.CheckAdmin(update)
		log.Println("admin")
	case "kick":

	default:
		fmt.Println("i dont't know this command")
		return
	}
}
