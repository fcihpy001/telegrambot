package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegramBot/group"
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

	case "filter":

	case "stop":

	case "filters":

	case "invite":
		group.GroupHandlerCommand(&update, bot.bot)
	case "stat":
		fallthrough
	case "stats":
		group.GroupHandlerCommand(&update, bot.bot)
	case "stat_week":
		group.GroupHandlerCommand(&update, bot.bot)
	case "mute":
		group.GroupHandlerCommand(&update, bot.bot)
	case "unmute":
		group.GroupHandlerCommand(&update, bot.bot)
	case "ban":
		group.GroupHandlerCommand(&update, bot.bot)
	case "unban":
		group.GroupHandlerCommand(&update, bot.bot)
	case "admin":
		group.GroupHandlerCommand(&update, bot.bot)
	case "kick":
		group.GroupHandlerCommand(&update, bot.bot)
	default:
		fmt.Println("i dont't know this command")
		return
	}
}
