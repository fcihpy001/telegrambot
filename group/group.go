package group

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegramBot/model"
	"telegramBot/utils"
)

func GroupHandlerQuery(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	query := update.CallbackQuery.Data
	switch query {
	case "group_setting":
		fmt.Println("group_setting")
	case "group_solitaire":
		fmt.Println("group_solitaire")
	case "group_record":
		fmt.Println("group_record")
	case "group_statistic":
		fmt.Println("group_statistic")
		mgr.statics(update)
	case "group_verification":
		fmt.Println("group_verification")
	case "group_welcome":
		mgr.welcomeNewMember(update.Message)
	case "group_speechtodayranging":
		mgr.speechRanging(update, "today")
	case "group_speech7daysranging":
		mgr.speechRanging(update, "week")
	case "group_speechstatistics":
		mgr.speechstatistics(update)
	case "group_invite_ranging":
		mgr.inviteRanging(update)
	case "group_invite_7days_ranging":
		mgr.invitestatis(update)
	case "group_today_quit":
		mgr.groupmemberstatis(update, "today")
	case "group_7days_quit":
		mgr.groupmemberstatis(update, "week")

	case "toast":
		fmt.Println("请选择")
	}
}

func (mgr *GroupManager) statics(update *tgbotapi.Update) {
	btn11 := model.ButtonInfo{
		Text:    ".发言统计.",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "今日发言排名",
		Data:    "group_speechtodayranging",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    "7日发言排名",
		Data:    "group_speech7daysranging",
		BtnType: model.BtnTypeData,
	}
	btn23 := model.ButtonInfo{
		Text:    "7日发言统计",
		Data:    "group_speechstatistics",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "📊邀请统计",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "今日邀请排名",
		Data:    "group_invite_ranging",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "7日邀请排名",
		Data:    "group_invite_7days_ranging",
		BtnType: model.BtnTypeData,
	}
	btn61 := model.ButtonInfo{
		Text:    "📊进退群统计",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn71 := model.ButtonInfo{
		Text:    "今日进退群数据",
		Data:    "group_today_quit",
		BtnType: model.BtnTypeData,
	}
	btn72 := model.ButtonInfo{
		Text:    "7日进退群统计",
		Data:    "group_7days_quit",
		BtnType: model.BtnTypeData,
	}
	btn81 := model.ButtonInfo{
		Text:    "返回首页",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btns := [][]model.ButtonInfo{{btn11}, {btn21, btn22, btn23}, {btn31}, {btn41, btn42}, {btn61}, {btn71, btn72}, {btn81}}
	keyboard := utils.MakeKeyboard(btns)
	utils.SendMenu(update.CallbackQuery.Message.Chat.ID,
		"📊 【流量聚集地】统计\n\n在群组中使用命令：\n/stat 查询今天活跃统计\n/stat_week 查询七天活跃统计\n/stats 自定义时间查询活跃统计\n\n查看命令帮助",
		keyboard,
		mgr.bot)

}

func (mgr *GroupManager) speechRanging(update *tgbotapi.Update, period string) {

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "今日总发言：7条，以下是前100名：\n\n1.fcihpy - 6 条\n2.Fcihpy3 - 1 条")
	if period == "week" {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "7日发言数：8条，以下是前100名：\n\n1.fcihpy - 7 条\n2.Fcihpy3 - 1 条")
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦀返回", "back"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *GroupManager) speechstatistics(update *tgbotapi.Update) {

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "以下是7日发言统计：\n\n2023-08-28       7 条\n2023-08-27       1 条")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦀返回", "back"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *GroupManager) inviteRanging(update *tgbotapi.Update) {

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "今日邀请：7人，以下是前100名：\n\n1.fcihpy - 6 条\n2.Fcihpy3 - 1 条")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦀返回", "back"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *GroupManager) invitestatis(update *tgbotapi.Update) {

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "7日邀请统计，以下是前100名：\n\n1.fcihpy - 6 条\n2.Fcihpy3 - 1 条")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦀返回", "back"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *GroupManager) groupmemberstatis(update *tgbotapi.Update, period string) {

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "今日进群：0人，退群：0人\n以下是今日最新进群20人：\n\n\n以下是今日最新退群20人：")
	if period == "week" {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "以下是过去7天进退群统计：\n入群共：0人，退群共：0人")
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦀返回", "back"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func GroupHandlerCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	query := update.CallbackQuery.Data

	switch query {
	case "invite":
		mgr.inviteLink(update)
	case "stats":

	case "stat_week":

	case "mute":

	case "unmute":

	case "ban":
		mgr.ban(update)
	case "unban":
		mgr.unBan(update)
	case "admin":
		mgr.checkAdmin(update)
	case "kick":

	default:
		fmt.Println("unknown command")
	}
}

func GroupHandlerMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	mgr.welcomeNewMember(message)
}
