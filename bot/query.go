package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/group"
	"telegramBot/lucky"
	"telegramBot/setting"
)

// 处理行内按钮事件
func (bot *SmartBot) handleQuery(update *tgbotapi.Update) {
	query := update.CallbackQuery.Data
	fmt.Println("query command--", query)
	if strings.HasPrefix(query, "lucky") {
		lucky.LuckyHandler(update, bot.bot)
	} else if strings.HasPrefix(query, "group") {
		group.GroupHandler(update, bot.bot)
	} else if strings.HasPrefix(query, "settings") {
		setting.Settings(update.CallbackQuery.Message.Chat.ID, bot.bot)
	} else {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "邀请链接")
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}

}
