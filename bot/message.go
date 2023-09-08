package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/group"
	"telegramBot/setting"
)

// 处理普通消息
func (bot *SmartBot) handleMessage(update *tgbotapi.Update) {
	// 获取用户发送的消息文本
	messageText := update.Message.Text
	fmt.Println("message handler:", messageText)

	if group.HandleAdminConversation(update, bot.bot) {
		// 如果消息被拦截 不需要后续处理
		return
	}

	if messageText == "webapp" {
		// 只能在私有聊天中使用
		group.SendTestWebApp(update, bot.bot)
		return
	}
	//违禁词处理
	setting.HandlerProhibited(update, bot.bot)

	//自动回复处理
	setting.HandlerAutoReply(update, bot.bot)

	return
	//处理自动回复

	if strings.HasPrefix(messageText, "统计") { //获取违禁词库
		reply := "今天活跃统计功能"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if strings.Contains(messageText, "美国") { //获取自动回复词库
		reply := "这是个违禁词，小心被禁言"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if strings.Contains(messageText, "reafw") { //获取活动关键词

	} else {
		reply := "感谢您的消息，我还在进修闭关中。。。"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
