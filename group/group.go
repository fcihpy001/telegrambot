package group

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GroupHandlerQuery(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	query := update.CallbackQuery.Data
	switch query {
	case "group_setting":
		fmt.Println("group_setting")
	case "group_solitaire":
		fmt.Println("group_solitaire")
	case "group_record":
		fmt.Println("group_record")
	case "group_statistic":
		fmt.Println("group_statistic")
	case "group_verification":
		fmt.Println("group_verification")
	case "group_welcome":
		mgr.welcomeNewMember(update.Message)
	}
}

func GroupHandlerCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	// if update.CallbackQuery == nil {
	// 	return
	// }
	query := strings.ToLower(update.Message.Command())
	switch query {
	case "invite":
		mgr.inviteLink(update)
	case "stats":
		mgr.StatsMemberMessages(update)

	case "stat_week":

	case "mute":
		mgr.Mute(update)

	case "unmute":
		mgr.UnMute(update)

	case "ban":
		mgr.ban(update)
	case "unban":
		mgr.unBan(update)
	case "admin":
		mgr.checkAdmin(update)
	case "kick":

	}
}

func GroupHandlerMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	mgr.welcomeNewMember(message)
}
