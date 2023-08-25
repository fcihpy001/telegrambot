package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func Help(chatId int64, bot *tgbotapi.BotAPI) {
	//TODO 获取当前群的名子
	msg := tgbotapi.NewMessage(chatId, "👏 欢迎使用ToplinkBot，如何使用：\n                \n •  邀请 @toplink 进入群组\n •  设置为管理员\n •  在机器人私聊中发送 /start 打开设置菜单。\n\n/help 查看我的功能\n\n\n👉 选择下面群组进行设置：")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("😺设置", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧠添加toplink到群组+", "group_join"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦴抽奖推送", "luckysetting"),
			tgbotapi.NewInlineKeyboardButtonData("👷‍订阅频道", "settings"),
			tgbotapi.NewInlineKeyboardButtonData("🎒官方群组", "settings"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
