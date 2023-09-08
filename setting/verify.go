package setting

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

	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	fmt.Println(query)
	if cmd == "verify_setting_menu" {
		verifySettingMenu(update, bot)

	} else if cmd == "verify_setting_status" {
		verifyStatusHandler(update, bot, params)

	} else if cmd == "verify_setting_method" {
		verifyMethodHandler(update, bot, params)

	} else if cmd == "verify_setting_time" {
		verifyTimeHandler(update, bot, params)

	} else if cmd == "verify_setting_punish" {
		verifyPunishHandler(update, bot, params)
	}
}

func verifySettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_ = services.GetModelData(utils.GroupInfo.GroupId, &verifySetting)
	verifySetting.ChatId = utils.GroupInfo.GroupId

	var buttons [][]model.ButtonInfo
	utils.Json2Button2("./config/verify.json", &buttons)
	fmt.Println(&buttons)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btnArr := buttons[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArr); j++ {
			btn := btnArr[j]
			updateVerifyButtonStatus(&btn)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.VerifySettingMenuMarkup = keyboard
	content := updateVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	bot.Send(msg)
}

func verifyStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
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

	content := updateVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	bot.Send(msg)

}

func verifyMethodHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
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
	content := updateVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	bot.Send(msg)

}
func verifyTimeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
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

	content := updateVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	bot.Send(msg)

}
func verifyPunishHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
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

	content := updateVerifySetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.VerifySettingMenuMarkup)
	bot.Send(msg)
}

func updateVerifySetting() string {
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

func updateVerifyButtonStatus(btn *model.ButtonInfo) {
	if btn.Text == "启用" && verifySetting.Enable {
		btn.Text = "✅启用"
	} else if btn.Text == "关闭" && !verifySetting.Enable {
		btn.Text = "✅关闭"
	} else if btn.Text == "按钮" && verifySetting.VerifyType == "按钮" {
		btn.Text = "✅按钮"
	} else if btn.Text == "数学题" && verifySetting.VerifyType == "数学题" {
		btn.Text = "✅数学题"
	} else if btn.Text == "验证码" && verifySetting.VerifyType == "验证码" {
		btn.Text = "✅验证码"
	} else if btn.Text == "1分" && verifySetting.VerifyTime == 1 {
		btn.Text = "✅1分"
	} else if btn.Text == "5分" && verifySetting.VerifyTime == 5 {
		btn.Text = "✅5分"
	} else if btn.Text == "10分" && verifySetting.VerifyTime == 10 {
		btn.Text = "✅10分"
	} else if btn.Text == "禁言" && verifySetting.PunishType == "禁言" {
		btn.Text = "✅禁言"
	} else if btn.Text == "踢出" && verifySetting.PunishType == "踢出" {
		btn.Text = "✅踢出"
	}
}
