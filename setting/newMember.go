package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var memberCheck model.NewMemberCheck

func MemberCheckHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	if cmd == "new_member_check_menu" {
		memberCheckMenu(update, bot)

	} else if cmd == "new_member_check_status" {
		memberCheckStatusHandler(update, bot, params == "enable")

	} else if cmd == "new_member_check_time_menu" {
		memberCheckTimeMenu(update, bot)
	}
}

func memberCheckMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	err := services.GetModelData(utils.GroupInfo.GroupId, &memberCheck)

	var buttons [][]model.ButtonInfo
	utils.Json2Button2("newMember.json", &buttons)
	fmt.Println(&buttons)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btnArr := buttons[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArr); j++ {
			row = append(row, btnArr[j])
		}
		rows = append(rows, row)
	}
	if memberCheck.Enable {
		rows[0][1].Text = "✅启用"
		rows[0][2].Text = "关闭"
	} else {
		rows[0][1].Text = "启用"
		rows[0][2].Text = "✅关闭"
	}
	keyboard := utils.MakeKeyboard(rows)
	utils.MemberCheckMarkup = keyboard

	//要读取用户设置的数据
	content := updateMemberSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func memberCheckStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
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

func memberCheckTimeMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

func MemberCheckTimeResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := "⚠️ 仅支持数字，请重新输入\n\n👉 请输入新群员限制时间（秒，例如：600）："
	if !isNumeric(update.Message.Text) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
		bot.Send(msg)
		return
	}
	content = "✅ 设置成功，点击按钮返回。"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "new_member_check_menu",
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

	memberCheck.ChatId = utils.GroupInfo.GroupId
	services.SaveModel(&memberCheck, utils.GroupInfo.GroupId)
	return content
}
