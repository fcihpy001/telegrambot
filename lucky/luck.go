package lucky

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// LuckyHandler 处理抽奖部分功能
func LuckyHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := LucyManager{
		bot: bot,
	}
	query := update.CallbackQuery.Data
	switch query {
	case "lucky_activity":
		mgr.luckyActivity(update)
	case "lucky_setting":
		mgr.luckysetting(update)
	case "lucky_create":
		mgr.luckyrecord(update)
	case "lucky_record":
		mgr.luckyrecord(update)
	}
}

func (mgr *LucyManager) luckyActivity(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🎁【测试】抽奖\n\n发起抽奖次数：0    \n\n已开奖：0       未开奖：0       取消：0")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📎发起抽奖活动", "lucky_create"),
			tgbotapi.NewInlineKeyboardButtonData("📪查看抽奖记录", "lucky_record"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧶设置抽奖", "lucky_setting"),
			tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *LucyManager) luckycreate(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🎁 测试发起抽奖\n  \n🔥 通用抽奖：群员在群内回复指定关键词参与抽奖\n\n🙋‍♂️ 指定群报道抽奖：A群成员进入B群回复指定关键词参与抽奖\n\n🪁 邀请人数抽奖：群成员用[专属链接]或[添加成员]拉人进群，到指定人数后参与抽奖\n\n🥰 群活跃抽奖：根据活跃排名抽奖，或达到活跃度参与随机抽奖\n\n🎰 娱乐抽奖：水果机、摇骰子、飞镖、保龄球....\n\n 选择抽奖类型：")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉抽奖抽奖", "createlucky"),
			tgbotapi.NewInlineKeyboardButtonData("🛎指定群报道抽奖", "luckyrecord"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📮邀请抽奖", "createlucky"),
			tgbotapi.NewInlineKeyboardButtonData("🏮群i", "luckyrecord"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁娱乐抽奖", "createlucky"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📡返回", "settings"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *LucyManager) luckyrecord(update *tgbotapi.Update) {

}

func (mgr *LucyManager) luckysetting(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "⚙ 抽奖设置\n\n✅ 发布置顶：\n└ 发布抽奖消息群内置顶\n✅ 结果置顶：\n└ 中奖结果消息群内置顶\n✅ 删除口令：\n└ 5分钟后自动删除群成员参加抽奖发送的口令消息")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉发布置顶", "createlucky"),
			tgbotapi.NewInlineKeyboardButtonData(" 启用", "luckyrecord"),
			tgbotapi.NewInlineKeyboardButtonData("✅关闭", "luckyrecord"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📮结果置顶", "createlucky"),
			tgbotapi.NewInlineKeyboardButtonData("启用", "luckyrecord"),
			tgbotapi.NewInlineKeyboardButtonData("✅关闭", "luckyrecord"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁娱乐抽奖", "createlucky"),
			tgbotapi.NewInlineKeyboardButtonData(" 启用", "luckyrecord"),
			tgbotapi.NewInlineKeyboardButtonData("✅关闭", "luckyrecord"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📡返回到抽奖", "settings"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
