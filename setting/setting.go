package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func Settings(chatId int64, bot *tgbotapi.BotAPI) {

	reply := "设置【测试】群组，选择要更改的项目"
	msg := tgbotapi.NewMessage(chatId, reply)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌺抽奖活动", "lucky_activity"),
			tgbotapi.NewInlineKeyboardButtonData("😊专属邀请链接生成", "invite_link"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍🎓群接龙", "group_solitaire"),
			tgbotapi.NewInlineKeyboardButtonData("🧝‍群统计", "group_statistic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🐞自动回复", "auto_reply"),
			tgbotapi.NewInlineKeyboardButtonData("🦊定时消息", "timing_message"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌳入群验证", "group_verification"),
			tgbotapi.NewInlineKeyboardButtonData("进群欢迎", "group_welcome"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦬反垃圾", "anti_spam"),
			tgbotapi.NewInlineKeyboardButtonData("🌓反刷屏", "anti_flood"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⛄️违禁词", "prohibited_words"),
			tgbotapi.NewInlineKeyboardButtonData("🌽用户检查", "user_check"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌗夜晚模式", "night_mode"),
			tgbotapi.NewInlineKeyboardButtonData("🌰新群员限制", "new_member_limit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚂下一页", "next_page"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥦语言切换", "language_switch"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏊切换其它群", "switch_group"),
			tgbotapi.NewInlineKeyboardButtonData("🪕打开群", "open_group"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
