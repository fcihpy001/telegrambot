package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// 处理普通消息
func (bot *SmartBot) handleMessage(update *tgbotapi.Update) {
	// 获取用户发送的消息文本
	messageText := update.Message.Text
	fmt.Println("message:", messageText)

	if strings.HasPrefix(messageText, "统计") {
		reply := "今天活跃统计功能"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else {
		reply := "感谢您的消息，我还在进修闭关中。。。"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
