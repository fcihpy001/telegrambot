package bot

import (
	"fmt"
	"strconv"
	"strings"
	"telegramBot/group"
	"telegramBot/services"
	"telegramBot/setting"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 处理以/开头的指令消息,如/help  /status等
func (bot *SmartBot) handleCommand(update tgbotapi.Update) {
	fmt.Println("command---", update.Message.Command())

	//处理页面跳转进来的
	if strings.HasPrefix(update.Message.Command(), "start") && update.Message.Chat.Type == "private" {
		//接收参数，取空格后面的内容
		args := strings.TrimSpace(strings.Replace(update.Message.Text, "/start", "", -1))
		if len(args) == 0 {
			setting.StartHandler(&update, bot.bot)
			return
		}
		//分割参数
		params := strings.Split(args, "_")

		//根据参数获取群组信息
		module := params[0]
		groupId, _ := strconv.Atoi(params[1])

		msg := update.Message
		if module == "manager" {
			where := fmt.Sprintf("group_id = %d and uid = %d", groupId, update.Message.From.ID)
			_ = services.GetModelWhere(where, &utils.GroupInfo)

			group.StartAdminConversation(int64(groupId),
				msg.Chat.ID,
				msg.From.ID,
				int64(msg.MessageID),
				msg.From.FirstName+" "+msg.From.LastName,
				group.ConversationStart,
				nil,
				nil,
			)
			//开始页面跳转
			setting.Settings(&update, bot.bot)
			return
		} else if module == "lucky" {
			group.StartAdminConversation(
				int64(groupId),
				msg.Chat.ID,
				msg.From.ID,
				int64(msg.MessageID),
				msg.From.FirstName+" "+msg.From.LastName,
				group.ConversationStart,
				nil,
				nil,
			)
			group.LuckyCreateIndex(&update, bot.bot,
				group.NewCallbackParam(int64(msg.Chat.ID), int(msg.From.ID), "", true),
			)
			return
		} else if module == "solitaire" {
			group.PlaySolitaire(&update, bot.bot, args)
			return
		}
		return
	}

	switch strings.ToLower(update.Message.Command()) {
	case "help":
		setting.Help(update.Message.Chat.ID, bot.bot)

	case "start":
		//判断是否是私聊
		if update.Message.Chat.Type == "private" {
			setting.StartHandler(&update, bot.bot)
		} else {
			//如果是管理员	弹出管理菜单
			member, _ := getMemberInfo(update.Message.Chat.ID, update.Message.From.ID, bot.bot)
			if member.IsAdministrator() || member.IsCreator() {
				setting.ManagerMenu(&update, bot.bot)
			}
		}
		println(update.Message.Text)
		// 如果参数中有solitaire: 开头 且在私有聊天中, 是用户接龙
		//args := strings.TrimSpace(strings.Replace(update.Message.Text, "/start", "", -1))
		//if strings.HasPrefix(args, "solitaire-") && update.Message != nil && update.Message.Chat.Type == "private" {
		//	group.PlaySolitaire(&update, bot.bot, args)
		//	return
		//}
		//setting.Settings(&update, bot.bot)

	case "filter":

	case "stop":

	case "filters":
		setting.ReplyCommandHandler(&update, bot.bot)

	case "link":
		setting.GetInviteLink(&update, bot.bot)

	case "stat", "stats", "statistic", "stat_week", "lucky", "create":
		group.GroupHandlerCommand(&update, bot.bot)

	case "mention":
		group.SendTestMentioned(bot.bot, &update)

	case "manager":
		setting.ManagerMenu(&update, bot.bot)

	case "ban", "kick", "unmute", "mute", "unban":
		setting.OperationHandler(&update, bot.bot)

	default:
		fmt.Println("i don't know this command")
		return
	}
}

func getMemberInfo(chat_id int64, user_id int64, bot *tgbotapi.BotAPI) (tgbotapi.ChatMember, error) {
	req := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chat_id,
			UserID: user_id,
		},
	}
	return bot.GetChatMember(req)
}
