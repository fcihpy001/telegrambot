package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegramBot/group"
	"telegramBot/model"
	"telegramBot/services"
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

	case "stat", "stats", "statistic", "stat_week", "mute", "unmute", "ban", "unban", "admin", "kick", "invite", "link":
		group.GroupHandlerCommand(&update, bot.bot)

	default:
		fmt.Println("i dont't know this command")
		return
	}
	//	todo: 保存用户信息
	u := model.User{
		Uid:  update.Message.From.ID,
		Name: update.Message.Chat.FirstName,
	}
	services.SaveUser(&u)
}
