package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/group"
	"telegramBot/lucky"
	"telegramBot/setting"
	"telegramBot/utils"
)

// 处理行内按钮事件
func (bot *SmartBot) handleQuery(update *tgbotapi.Update) {
	query := update.CallbackQuery.Data
	fmt.Println("query command--", query)

	if strings.HasPrefix(query, "lucky") {
		lucky.LuckyHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "group") {
		group.GroupHandlerQuery(update, bot.bot)

	} else if strings.HasPrefix(query, "settings") {
		setting.Settings(update.CallbackQuery.Message.Chat.ID, bot.bot)

	} else if query == "join_group" {
		fmt.Println("replay...")
		// 创建 ForceReply 结构
		forceReply := tgbotapi.ForceReply{
			ForceReply: true,
		}

		// 创建包含 ForceReply 的消息
		message := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "请回复此消息：")
		message.ReplyMarkup = forceReply

		// 发送消息
		_, err := bot.bot.Send(message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	} else if query == "next_page" {
		//	发送还键盘的推送消息
		msg := tgbotapi.NewMessage(6401399435, "🎁【零度社区 (LingduDAO)- 中文群】群组发起了发言次数抽奖活动\n已开奖：1       未开奖：1       取消：0\n\nLDD是零度DAO的社区币\n├参与条件：发言6条\n├发言起始统计时间：2023-08-28 11:20:00\n├开奖时间：2023-08-28 22:00:00\n├奖品列表：\n├       2USDT     ×3份\n\n【如何参与？】在群组中发言6次，参与活动。")
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🌺加入活动群众", "lucky_activity"),
			))
		msg.ReplyMarkup = inlineKeyboard
		msg.DisableNotification = false
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if query == "prohibited_words" { //违禁词处理
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🔇 违禁词\n\n👉请输入添加的违禁词（一行一个）：")
		replayKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("返回"),
			))
		msg.ReplyMarkup = replayKeyboard
		msg.DisableNotification = false
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if query == "go_setting" { //违禁词列表
		bot.go_setting(update)
	} else {
		msg := tgbotapi.NewMessage(6401399435, "测试推送事件")
		msg.DisableNotification = false
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
	utils.SendReply(update.CallbackQuery.ID, bot.bot, false, "消息已经处理")
}

func (bot *SmartBot) go_setting(update *tgbotapi.Update) {
	fmt.Println("go_setting...")
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "设置【流量聚集地】群组，选择要更改的项目", utils.SettingMenuMarkup)
	bot.bot.Send(msg)
}
