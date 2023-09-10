package group

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func GroupHandlerQuery(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	query := update.CallbackQuery.Data
	fmt.Println("group query --", query)
	switch query {

	case "group_solitaire":
		fmt.Println("group_solitaire")
		mgr.SolitaireIndex(update)

	case "group_record":
		fmt.Println("group_record")
	case "group_statistic":
		fmt.Println("group_statistic")
		mgr.statics(update)

	case "group_speechtodayranging":
		mgr.speechRanging(update, "today")
	case "group_speech7daysranging":
		mgr.speechRanging(update, "week")
	case "group_speechstatistics":
		mgr.speechstatistics(update)
	case "group_invite_ranging":
		mgr.invitesToday(update)
	case "group_invite_7days_ranging":
		mgr.invitesWeek(update)
	case "group_today_quit":
		mgr.groupmemberstatis(update, "today")
	case "group_7days_quit":
		mgr.groupmemberstatis(update, "week")
	case "group_back_statics":
		mgr.group_back_statics(update)

	case "toast":
		fmt.Println("请选择")
	}
}

func GroupHandlerCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}

	query := strings.ToLower(update.Message.Command())
	switch query {
	case "invite":
		//mgr.inviteLink(update)
	case "stats":
		mgr.StatsMemberMessages(update)
	case "stat_week":

	case "stat":
		mgr.StatsMemberMessages(update)
	case "mute":
		mgr.UnMute(update)
	case "unmute":
		mgr.UnMute(update)
	case "ban":
		mgr.ban(update)
	case "unban":
		mgr.unBan(update)
	case "admin":
		mgr.checkAdmin(update)
	case "kick":

	default:
		fmt.Println("unknown command")
	}
}

func GroupHandlerMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	if message.NewChatMembers != nil {
		mgr.welcomeNewMember(message)
		return
	}
}
