package group

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var verifySetting model.VerifySetting

// 入群验证
func VerifySettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_ = services.GetModelData(update.CallbackQuery.Message.Chat.ID, &verifySetting)
	verifySetting.ChatId = update.CallbackQuery.Message.Chat.ID

	mgr := GroupManager{
		bot: bot,
	}
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	fmt.Println(query)
	if cmd == "verify_setting" {
		mgr.verifySettingMenu(update)

	} else if cmd == "verify_setting_status" {
		mgr.verifyStatusHandler(update, params)

	} else if cmd == "verify_setting_method" {
		mgr.verifyMethodHandler(update, params)

	} else if cmd == "verify_setting_time" {
		mgr.verifyTimeHandler(update, params)

	} else if cmd == "verify_setting_punish" {
		mgr.verifyPunishHandler(update, params)

	}

}

func (mgr *GroupManager) verifySettingMenu(update *tgbotapi.Update) {

	btn12txt := "启用"
	btn13txt := "✅关闭"
	if verifySetting.Enable {
		btn12txt = "✅启用"
		btn13txt = "关闭"
	}

	btn22txt := "按钮"
	btn23txt := "✅数学题"
	btn24txt := "验证码"
	if verifySetting.VerifyType == "按钮" {
		btn22txt = "✅按钮"
		btn23txt = "数学题"
		btn24txt = "验证码"
	} else if verifySetting.VerifyType == "数学题" {
		btn22txt = "按钮"
		btn23txt = "✅数学题"
		btn24txt = "验证码"
	} else if verifySetting.VerifyType == "验证码" {
		btn22txt = "按钮"
		btn23txt = "数学题"
		btn24txt = "✅验证码"
	}

	btn32txt := "1分"
	btn33txt := "5分"
	btn34txt := "10分"
	if verifySetting.VerifyTime == 1 {
		btn32txt = "✅1分"
		btn33txt = "5分"
		btn34txt = "10分"
	} else if verifySetting.VerifyTime == 5 {
		btn32txt = "1分"
		btn33txt = "✅5分"
		btn34txt = "10分"
	} else if verifySetting.VerifyTime == 10 {
		btn32txt = "1分"
		btn33txt = "5分"
		btn34txt = "✅10分"
	}

	btn42txt := "禁言"
	btn43txt := "✅踢出"
	if verifySetting.PunishType == "禁言" {
		btn42txt = "✅禁言"
		btn43txt = "踢出"
	} else if verifySetting.PunishType == "踢出" {
		btn42txt = "禁言"
		btn43txt = "✅踢出"
	}

	btn11 := model.ButtonInfo{
		Text:    "是否启用",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "verify_setting_status:enable",
		BtnType: model.BtnTypeData,
	}
	btn13 := model.ButtonInfo{
		Text:    btn13txt,
		Data:    "verify_setting_status:disable",
		BtnType: model.BtnTypeData,
	}

	btn21 := model.ButtonInfo{
		Text:    "模式",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    btn22txt,
		Data:    "verify_setting_method:按钮",
		BtnType: model.BtnTypeData,
	}
	btn23 := model.ButtonInfo{
		Text:    btn23txt,
		Data:    "verify_setting_method:数学题",
		BtnType: model.BtnTypeData,
	}
	btn24 := model.ButtonInfo{
		Text:    btn24txt,
		Data:    "verify_setting_method:验证码",
		BtnType: model.BtnTypeData,
	}

	btn31 := model.ButtonInfo{
		Text:    "验证时间",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn32 := model.ButtonInfo{
		Text:    btn32txt,
		Data:    "verify_setting_time:1",
		BtnType: model.BtnTypeData,
	}
	btn33 := model.ButtonInfo{
		Text:    btn33txt,
		Data:    "verify_setting_time:5",
		BtnType: model.BtnTypeData,
	}

	btn34 := model.ButtonInfo{
		Text:    btn34txt,
		Data:    "verify_setting_time:10",
		BtnType: model.BtnTypeData,
	}

	btn41 := model.ButtonInfo{
		Text:    "超时处理",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    btn42txt,
		Data:    "verify_setting_punish:禁言",
		BtnType: model.BtnTypeData,
	}
	btn43 := model.ButtonInfo{
		Text:    btn43txt,
		Data:    "verify_setting_punish:踢出",
		BtnType: model.BtnTypeData,
	}

	btn51 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}

	row1 := []model.ButtonInfo{btn11, btn12, btn13}
	row2 := []model.ButtonInfo{btn21, btn22, btn23, btn24}
	row3 := []model.ButtonInfo{btn31, btn32, btn33, btn34}
	row4 := []model.ButtonInfo{btn41, btn42, btn43}
	row5 := []model.ButtonInfo{btn51}

	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5}
	keyboard := utils.MakeKeyboard(rows)
	utils.VerifySettingMenuMarkup = keyboard
	content := updateGroupVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	mgr.bot.Send(msg)
}

func (mgr *GroupManager) verifyStatusHandler(update *tgbotapi.Update, params string) {
	if len(params) == 0 {
		return
	}
	if params == "enable" {
		verifySetting.Enable = true
		utils.VerifySettingMenuMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.VerifySettingMenuMarkup.InlineKeyboard[0][2].Text = "关闭"
	} else {
		verifySetting.Enable = false
		utils.VerifySettingMenuMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.VerifySettingMenuMarkup.InlineKeyboard[0][2].Text = "✅关闭"
	}

	content := updateGroupVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	mgr.bot.Send(msg)

}

func (mgr *GroupManager) verifyMethodHandler(update *tgbotapi.Update, params string) {
	if len(params) == 0 {
		return
	}

	verifySetting.VerifyType = params
	if params == "按钮" {
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][1].Text = "✅按钮"
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][2].Text = "数学题"
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][3].Text = "验证码"

	} else if params == "数学题" {
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][1].Text = "按钮"
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][2].Text = "✅数学题"
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][3].Text = "验证码"

	} else if params == "验证码" {
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][1].Text = "按钮"
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][2].Text = "数学题"
		utils.VerifySettingMenuMarkup.InlineKeyboard[1][3].Text = "✅验证码"

	}

	content := updateGroupVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	mgr.bot.Send(msg)

}
func (mgr *GroupManager) verifyTimeHandler(update *tgbotapi.Update, params string) {
	if len(params) == 0 {
		return
	}
	if params == "1" {
		verifySetting.VerifyTime = 1
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][1].Text = "✅1分"
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][2].Text = "5分"
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][3].Text = "10分"

	} else if params == "5" {
		verifySetting.VerifyTime = 5
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][1].Text = "1分"
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][2].Text = "✅5分"
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][3].Text = "10分"

	} else if params == "10" {
		verifySetting.VerifyTime = 10
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][1].Text = "1分"
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][2].Text = "5分"
		utils.VerifySettingMenuMarkup.InlineKeyboard[2][3].Text = "✅10分"
	}

	content := updateGroupVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	mgr.bot.Send(msg)

}
func (mgr *GroupManager) verifyPunishHandler(update *tgbotapi.Update, params string) {
	if len(params) == 0 {
		return
	}
	verifySetting.PunishType = params
	if params == "禁言" {
		utils.VerifySettingMenuMarkup.InlineKeyboard[3][1].Text = "✅禁言"
		utils.VerifySettingMenuMarkup.InlineKeyboard[3][2].Text = "踢出"
	} else if params == "踢出" {
		utils.VerifySettingMenuMarkup.InlineKeyboard[3][1].Text = "禁言"
		utils.VerifySettingMenuMarkup.InlineKeyboard[3][2].Text = "✅踢出"
	}

	content := updateGroupVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	mgr.bot.Send(msg)

}

func updateGroupVerifySetting() string {
	content := "🤖 入群验证\n启用后，用户进入群组需要验证才能发送消息\n\n"
	if verifySetting.Enable {
		content += "当前状态：✅已启用\n"
	} else {
		content += "当前状态：❌关闭\n"
	}
	content += "验证方法：" + verifySetting.VerifyType + "\n"

	content += "验证时间：" + fmt.Sprintf("%d", verifySetting.VerifyTime) + "分钟\n"

	content += "超时处理: " + verifySetting.PunishType + "\n"
	services.SaveModel(&verifySetting, verifySetting.ChatId)
	return content
}
