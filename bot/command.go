package bot

import (
	"fmt"
	"log"
	"strings"
	"telegramBot/group"
	"telegramBot/setting"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 处理以/开头的指令消息,如/help  /status等
func (bot *SmartBot) handleCommand(update tgbotapi.Update) {
	fmt.Println("command---", update.Message.Command())

	switch strings.ToLower(update.Message.Command()) {
	case "help":
		setting.Help(update.Message.Chat.ID, bot.bot)

	case "start":
		setting.StartHandler(&update, bot.bot)
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

	case "stat", "stats", "statistic", "stat_week", "mute", "unmute", "ban", "unban", "admin", "kick", "invite", "link":
		group.GroupHandlerCommand(&update, bot.bot)

	case "mention":
		group.SendTestMentioned(bot.bot, &update)

	case "test":
		// 创建投票参数
		poll := tgbotapi.NewPoll(update.Message.Chat.ID, "toplink 发起投票", "选项1", "选项2", "选项3")

		// 添加投票选项
		poll.Question = "这是一个投票qfqgq"

		poll.Options = []string{"选项11", "选项12", "选项13"}

		poll.AllowSendingWithoutReply = true // 允许在没有回复的情况下发送投票
		poll.ChannelUsername = "toplink"     // 投票频道的用户名

		// 设置其他投票参数
		poll.IsAnonymous = false           // 是否匿名投票
		poll.AllowsMultipleAnswers = false // 是否允许多选
		poll.OpenPeriod = 30               // 投票持续时间（以秒为单位）

		// 发送投票
		_, err := bot.bot.Send(poll)
		if err != nil {
			log.Println(err)
		}
	default:
		fmt.Println("i dont't know this command")
		return
	}
}
