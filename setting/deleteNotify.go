package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/utils"
)

// 初始选中信息
var deleteNotifySelect model.SelectInfo = model.SelectInfo{
	Row:    0,
	Column: 0,
	Text:   "10秒",
}
var notifyClass string
var deleteTime = 10

func DeleteNotifyHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	if cmd == "delete_notify_menu" { //删除提醒消息主菜单
		notifyMenu(update, bot, params)

	} else if cmd == "delete_notify_time" { //删除提醒消息时间设置
		notifyTimeHandler(update, bot, params)
	}
}

func notifyMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	notifyClass = params

	times := utils.GetTimeData()

	var rows [][]model.ButtonInfo
	for i := 0; i < len(times); i++ {
		timeArray := times[i]
		var row []model.ButtonInfo
		for j := 0; j < len(timeArray); j++ {
			time := times[i][j]
			btn := model.ButtonInfo{
				Text:    time,
				Data:    "delete_notify_time:" + strconv.Itoa(i) + "&" + strconv.Itoa(j) + "&" + time,
				BtnType: model.BtnTypeData,
			}
			row = append(row, btn)
		}
		rows = append(rows, row)
	}
	//返回按钮
	backBtn := model.ButtonInfo{
		Text:    "返回",
		Data:    "prohibited_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row := []model.ButtonInfo{backBtn}
	rows = append(rows, row)

	keyboard := utils.MakeKeyboard(rows)
	utils.DeleteNotifyMenuMarkup = keyboard

	content := updateDeleteNotifyMsg()
	sendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)
}

func notifyTimeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	data := strings.Split(params, "&")
	row, _ := strconv.Atoi(data[0])
	column, _ := strconv.Atoi(data[1])
	text := data[2]

	deleteTime = utils.ParseTime(text)

	//取消以前的选中
	utils.DeleteNotifyMenuMarkup.InlineKeyboard[deleteNotifySelect.Row][deleteNotifySelect.Column].Text = deleteNotifySelect.Text
	//更新选中
	utils.DeleteNotifyMenuMarkup.InlineKeyboard[row][column].Text = "✅" + text
	//更新选中信息
	deleteNotifySelect.Row = row
	deleteNotifySelect.Column = column
	deleteNotifySelect.Text = text

	content := updateDeleteNotifyMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.DeleteNotifyMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

func updateDeleteNotifyMsg() string {
	content := ""
	if notifyClass == "prohibited" {
		content = "🔇 违禁词\n\n群成员触发🔇 违禁词时，机器人发出的提醒消息在多少时间后自动删除"
		prohibitedSetting.DeleteNotifyMsgTime = deleteTime
		updateProhibitedSettingMsg()

	} else if notifyClass == "flood" {
		content = "💬 反刷屏\n\n群成员触发💬 反刷屏时，机器人发出的提醒消息在多少时间后自动删除"
		floodSetting.DeleteNotifyMsgTime = deleteTime
		updateFloodMsg()

	} else if notifyClass == "spam" {
		content = "📨 反垃圾\n\n群成员触发📨 反垃圾时，机器人发出的提醒消息在多少时间后自动删除"
		spamsSetting.DeleteNotifyMsgTime = deleteTime
		updateSpamMsg()

	} else if notifyClass == "userCheck" {
		content = "🔦 用户检查\n\n群成员触发🔦 用户检查时，机器人发出的提醒消息在多少时间后自动删除"
		userCheckSetting.DeleteNotifyMsgTime = deleteTime
		updateUserSettingMsg()
	}
	return content
}
