package bot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"telegramBot/group"
	"telegramBot/services"
	"telegramBot/setting"
	"telegramBot/utils"
)

// 处理以/开头的指令消息,如/help  /status等
func (bot *SmartBot) handleCommand(update tgbotapi.Update) {
	fmt.Println("command---", update.Message.Command())

	//处理页面跳转进来的
	if strings.HasPrefix(update.Message.Command(), "start") && update.Message.Chat.Type == "private" {
		//接收参数，取空格后面的内容
		args := strings.TrimSpace(strings.Replace(update.Message.Text, "/start", "", -1))
		fmt.Println("args", args)

		if len(args) > 0 {
			//根据参数获取群组信息
			groupId, _ := strconv.Atoi(args)
			where := fmt.Sprintf("group_id = %d and uid = %d", groupId, update.Message.From.ID)
			_ = services.GetModelWhere(where, &utils.GroupInfo)

			//开始页面跳转
			setting.Settings(&update, bot.bot)
		} else {
			setting.StartHandler(&update, bot.bot)
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

	case "setting":
		// 这里如果有参数, 进入对应的处理逻辑; 否则展示管理界面
		println(update.Message.Text)
		// 如果参数中有solitaire: 开头 且在私有聊天中, 是用户接龙
		args := strings.TrimSpace(strings.Replace(update.Message.Text, "/start", "", -1))
		if strings.HasPrefix(args, "solitaire-") && update.Message != nil && update.Message.Chat.Type == "private" {
			group.PlaySolitaire(&update, bot.bot, args)
			return
		}
		setting.Settings(&update, bot.bot)

	case "create":

	case "luck":

	case "filter":

	case "stop":

	case "filters":

	case "link":
		setting.GetInviteLink(&update, bot.bot)

	case "stat", "stats", "statistic", "stat_week", "admin", "invite":
		group.GroupHandlerCommand(&update, bot.bot)

	case "mention":
		group.SendTestMentioned(bot.bot, &update)

	case "manager":
		setting.ManagerMenu(&update, bot.bot)

	case "ban", "kick", "unmute", "mute", "unban":
		setting.OperationHandler(&update, bot.bot)

	case "test":
		//setting.ScheduleMessage(&update, bot.bot)
		testapp(bot.bot, "https://python-telegram-bot.org/static/webappbot", "点击这里了解信息", update.Message.Chat.ID, "这是个好东西")

	default:
		fmt.Println("i dont't know this command")
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

func testapp(bot *tgbotapi.BotAPI, url string, buttonTitle string, receiver int64, desc string) {
	data := make(map[string]interface{})
	data["inline_keyboard"] = [][]interface{}{
		{
			map[string]interface{}{
				"text": buttonTitle,
				"web_app": map[string]string{
					"url": url,
				},
			},
		},
	}
	payload, _ := json.Marshal(data)

	params := map[string]string{
		"chat_id":      fmt.Sprint(receiver),
		"text":         desc,
		"reply_markup": string(payload), //
	}

	resp, err := bot.MakeRequest("sendMessage", params)
	if err != nil {
		log.Println(err)
	}
	buf, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(buf))
	tgbotapi.NewMessage(receiver, "ok")
}
