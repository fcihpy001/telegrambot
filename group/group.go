package group

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
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
		mgr.statics(update)
	case "group_verification":
		fmt.Println("group_verification")
	case "group_welcome":
		mgr.welcomeNewMember(update.Message)
	case "group_speechtodayranging":
		mgr.speechRanging(update, "today")
	case "group_speech7daysranging":
		mgr.speechRanging(update, "week")
	case "group_speechstatistics":
		mgr.speechstatistics(update)
	case "group_invite_ranging":
		mgr.inviteRanging(update)
	case "group_invite_7days_ranging":
		mgr.invitestatis(update)
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
	fmt.Println("groupcommand---", update.Message.Command())

	switch strings.ToLower(update.Message.Command()) {
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

	case "link":
		mgr.getInviteLink(update.Message.Chat.ID, update.Message.From.FirstName)

	default:
		fmt.Println("unknown command")
	}
}

func GroupHandlerMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	mgr.welcomeNewMember(message)
}

func (mgr *GroupManager) getInviteLink(receiver int64, name string) {
	config := tgbotapi.CreateChatInviteLinkConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: receiver,
		},
		Name:               "fcihpy",
		ExpireDate:         int(time.Now().Unix() + 86400*365),
		MemberLimit:        9999,
		CreatesJoinRequest: false,
	}
	resp, err := mgr.bot.Request(config)
	if err != nil {
		fmt.Println("linkerr111", err)
	}
	m := map[string]interface{}{}
	json.Unmarshal(resp.Result, &m)
	link := m["invite_link"].(string)

	msg := fmt.Sprintf("🔗 %s 您的专属链接:\n %s (点击复制)\n\n👉 👉 当前总共邀请0人\n\n（本消息5分钟自毁）", name, link)
	mgr.sendText(receiver, msg)
}
