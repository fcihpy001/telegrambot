package group

import (
	"errors"
	"log"
	"telegramBot/model"
	"telegramBot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// LuckyHandler 处理抽奖部分功能
func LuckyHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
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

func (mgr *GroupManager) luckyActivity(update *tgbotapi.Update) {
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

func (mgr *GroupManager) luckycreate(update *tgbotapi.Update) {
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

func (mgr *GroupManager) luckyrecord(update *tgbotapi.Update) {

}

func (mgr *GroupManager) luckysetting(update *tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatId, "⚙ 抽奖设置\n\n✅ 发布置顶：\n└ 发布抽奖消息群内置顶\n✅ 结果置顶：\n└ 中奖结果消息群内置顶\n✅ 删除口令：\n└ 5分钟后自动删除群成员参加抽奖发送的口令消息")

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉发布置顶", ""),
			tgbotapi.NewInlineKeyboardButtonData(" 启用", "luckysetting?typ=startPin&val=1"),
			tgbotapi.NewInlineKeyboardButtonData("✅关闭", "luckysetting?typ=startPin&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📮结果置顶", ""),
			tgbotapi.NewInlineKeyboardButtonData("启用", "luckysetting?typ=endPin&val=1"),
			tgbotapi.NewInlineKeyboardButtonData("✅关闭", "luckysetting?typ=endPin&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁删除口令", ""),
			tgbotapi.NewInlineKeyboardButtonData(" 启用", "luckysetting?typ=deleteToken&val=1"),
			tgbotapi.NewInlineKeyboardButtonData("✅关闭", "luckysetting?typ=deleteToken&val=0"),
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

func luckyIndex(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📎发起抽奖活动", "lucky_create"),
			tgbotapi.NewInlineKeyboardButtonData("📪查看抽奖记录", "lucky_record"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧶设置抽奖", "luckysetting"),
			tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
		))
	// todo
	msg := tgbotapi.NewEditMessageTextAndMarkup(param.chatId, param.msgId,
		"🎁【测试】抽奖\n\n发起抽奖次数：0    \n\n已开奖：0       未开奖：0       取消：0", inlineKeyboard)

	_, err := bot.Send(msg)
	if err != nil {
		logger.Err(err).Msg("send lucky index failed")
	}
	return err
}

func toggleLuckySetting(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	var (
		toggleType string
		toggleVal  bool
	)

	logger.Info().Str("param", param.data).Msg("toggle lucky setting param")

	chatSetting := services.FindChatLuckySetting(param.chatId)
	if chatSetting == nil {
		// 初始值
		chatSetting = &model.LuckySetting{
			ChatId:       param.chatId,
			StartPinned:  true,
			ResultPinned: true,
			DeleteToken:  true,
		}
	}

	if len(param.param) > 0 {
		typ := param.param["typ"]
		if len(typ) == 0 {
			err := errors.New("toggleLuckySetting: not found param toggle type")
			logger.Error().Msg("toggleLuckySetting: not found param toggle type")
			return err
		}
		toggleType = typ[0]
		val := param.param["val"]
		if len(val) == 0 {
			err := errors.New("toggleLuckySetting: not found param toggle value")
			logger.Error().Msg("toggleLuckySetting: not found param toggle value")
			return err
		}
		toggleVal = toBool(val[0])
		switch toggleType {
		case "startPin":
			chatSetting.StartPinned = toggleVal
		case "endPin":
			chatSetting.ResultPinned = toggleVal
		case "deleteToken":
			chatSetting.DeleteToken = toggleVal
		}
		// update chat lucky setting
		services.UpdateChatLuckySettings(chatSetting)
	}

	content := "⚙ 抽奖设置\n\n✅ 发布置顶：\n└ 发布抽奖消息群内置顶\n✅ 结果置顶：\n└ 中奖结果消息群内置顶\n✅ 删除口令：\n└ 5分钟后自动删除群成员参加抽奖发送的口令消息"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉发布置顶", "luckysetting?alert=xxx"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(chatSetting.StartPinned, "启用"),
				"luckysetting?typ=startPin&val=1"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(!chatSetting.StartPinned, "关闭"),
				"luckysetting?typ=startPin&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📮结果置顶", "luckysetting?alert=xxx"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(chatSetting.ResultPinned, "启用"),
				"luckysetting?typ=endPin&val=1"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(!chatSetting.ResultPinned, "关闭"),
				"luckysetting?typ=endPin&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁删除口令", "luckysetting?alert=xxx"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(chatSetting.DeleteToken, "启用"),
				"luckysetting?typ=deleteToken&val=1"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(!chatSetting.DeleteToken, "关闭"),
				"luckysetting?typ=deleteToken&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📡返回到抽奖", "lucky$"),
		))
	reply := tgbotapi.NewEditMessageTextAndMarkup(param.chatId, param.msgId, content, inlineKeyboard)

	_, err := bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send toggleLuckySetting message failed")
	}
	return err
}
