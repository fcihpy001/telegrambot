package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegramBot/setting"
)

// 处理以/开头的指令消息,如/help  /status等
func (bot *SmartBot) handleReply(update *tgbotapi.Update) {
	fmt.Println("reply---", update.Message.ReplyToMessage.Text)
	replyMsg := update.Message.ReplyToMessage.Text
	if strings.Contains(replyMsg, "输入添加的违禁词（一行一个") {
		setting.ProhibitedAdd(update, bot.bot)
	}
}
