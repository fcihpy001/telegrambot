package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var memberCheck model.NewMemberCheck

func MemberCheckMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		memberCheck = services.GetMemberSettings(update.CallbackQuery.Message.Chat.ID)
	}

	memberCheck.ChatId = update.CallbackQuery.Message.Chat.ID

	btn12txt := "启用"
	btn13txt := "✅关闭"
	if memberCheck.Enable {
		btn12txt = "✅启用"
		btn13txt = "关闭"
	}

	btn11 := model.ButtonInfo{
		Text:    "限制发消息",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "new_member_check_enable",
		BtnType: model.BtnTypeData,
	}
	btn13 := model.ButtonInfo{
		Text:    btn13txt,
		Data:    "new_member_check_disable",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "设置限制时间",
		Data:    "new_member_check_time_menu",
		BtnType: model.BtnTypeData,
	}

	btn31 := model.ButtonInfo{
		Text:    "🏠返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12, btn13}
	row2 := []model.ButtonInfo{btn21}
	row3 := []model.ButtonInfo{btn31}
	rows := [][]model.ButtonInfo{row1, row2, row3}
	keyboard := utils.MakeKeyboard(rows)
	utils.MemberCheckMarkup = keyboard

	//要读取用户设置的数据
	content := updateMemberSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func MemberCheckStatus(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
	if enable {
		utils.MemberCheckMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.MemberCheckMarkup.InlineKeyboard[0][2].Text = "关闭"
	} else {
		utils.MemberCheckMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.MemberCheckMarkup.InlineKeyboard[0][2].Text = "✅关闭"
	}
	memberCheck.Enable = enable

	content := updateMemberSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.MemberCheckMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func MemberCheckTimeMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	time := ""
	if memberCheck.DelayTime < 61 {
		time = fmt.Sprintf("%d秒", memberCheck.DelayTime)
	} else if memberCheck.DelayTime > 60 {
		time = fmt.Sprintf("%d分钟", memberCheck.DelayTime/60)
	}
	content := fmt.Sprintf("👤 新群员限制\n\n当前设置：%s\n👉 请输入新群员限制时间（秒，例如：600）：", time)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	keybord := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("返回"),
		))

	msg.ReplyMarkup = keybord
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}

	bot.Send(msg)
}

func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func MemberCheckTimeAction(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := "⚠️ 仅支持数字，请重新输入\n\n👉 请输入新群员限制时间（秒，例如：600）："
	if !isNumeric(update.Message.Text) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
		bot.Send(msg)
		return
	}
	content = "✅ 设置成功，点击按钮返回。"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "new_member_check",
		BtnType: model.BtnTypeData,
	}
	memberCheck.DelayTime, _ = strconv.Atoi(update.Message.Text)

	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)
	updateMemberSettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateMemberSettingMsg() string {
	content := "👤 新群员限制\n\n"
	enableMsg := "❌限制发消息：\n"
	if memberCheck.Enable {
		enableMsg = "✅限制发消息：\n"
	}
	time := ""
	if memberCheck.DelayTime < 61 {
		time = fmt.Sprintf("%d秒", memberCheck.DelayTime)
	} else if memberCheck.DelayTime > 60 {
		time = fmt.Sprintf("%d分钟", memberCheck.DelayTime/60)
	}
	limitTime := fmt.Sprintf("└ 新群员进群在设置时间 %s 内，不能发送消息", time)

	content = content + enableMsg + limitTime
	services.SaveMemberSettings(&memberCheck)
	return content
}
