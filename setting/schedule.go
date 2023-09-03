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

var scheduleMsg = model.ScheduleMsg{}
var selectInfo = model.SelectInfo{}

func ScheduleSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_ = services.GetModelData(update.CallbackQuery.Message.Chat.ID, &scheduleMsg)
	scheduleMsg.ChatId = update.CallbackQuery.Message.Chat.ID

	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	fmt.Println(query)
	if cmd == "schedule_message" { //定时消息设置主菜单
		scheduleMenu(update, bot)

	} else if cmd == "schedule_and" { //添加定时消息
		scheduleMsgMenu(update, bot)

	} else if cmd == "schedule_status" { //定时消息状态设置
		scheduleStatusHandler(update, bot, params)

	} else if cmd == "schedule_delete_prev" { //定时消息删除上一条消息设置
		scheduleDeletePrevHandler(update, bot, params)

	} else if cmd == "schedule_pin" { //定时消息置顶设置
		schedulePinHandler(update, bot, params)

	} else if cmd == "schedule_repeat_hour" {
		repeatValueHandler(update, bot, params)

	} else if cmd == "schedule_repeat_minute" {
		repeatValueHandler(update, bot, params)

	} else if cmd == "schedule_repeat" { //定时消息重复频率设置
		scheduleRepeatMenuHandler(update, bot)

	} else if cmd == "schedule_time" { //定时消息时间段设置
		scheduleTimeStartMenuHandler(update, bot)

	} else if strings.HasPrefix(cmd, "schedule_time") { //定时消息时间段内容设置
		scheduleTimeChangeHandler(update, bot)

	} else if cmd == "schedule_date_start" {
		scheduleDateStartHandler(update, bot)

	} else if cmd == "schedule_date_end" {
		scheduleDateEndHandler(update, bot)

	} else if cmd == "schedule_content" {
		scheduleContentHandler(update, bot)

	}
}

func repeatValueHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	selectInfo.Text = params

	msg1 := "✅启用"
	msg2 := "关闭"
	scheduleMsg.Enable = true
	if params == "disable" {
		scheduleMsg.Enable = false
		msg1 = "启用"
		msg2 = "✅关闭"
	}
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[0][1].Text = msg1
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[0][2].Text = msg2
	content := updateScheduleMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ScheduleMsgMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

func scheduleMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	btnAnd := model.ButtonInfo{
		Text:    "添加定时消息",
		Data:    "schedule_and",
		BtnType: model.BtnTypeData,
	}
	btnBack := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btnAnd}
	row2 := []model.ButtonInfo{btnBack}
	rows := [][]model.ButtonInfo{row1, row2}
	keyboard := utils.MakeKeyboard(rows)
	content := updateScheduleList()
	utils.SendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)

}

func scheduleMsgMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	btn01txt := "启用"
	btn02txt := "✅关闭"
	if scheduleMsg.Enable {
		btn01txt = "✅启用"
		btn02txt = "关闭"
	}
	btn11txt := "✅是"
	btn12txt := "否"
	if !scheduleMsg.DeletePrevMsg {
		btn11txt = "是"
		btn12txt = "✅否"
	}
	btn21txt := "✅是"
	btn22txt := "否"
	if scheduleMsg.Pin {
		btn21txt = "是"
		btn22txt = "✅否"
	}

	btn00 := model.ButtonInfo{
		Text:    "状态",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn01 := model.ButtonInfo{
		Text:    btn01txt,
		Data:    "schedule_status:enable",
		BtnType: model.BtnTypeData,
	}
	btn02 := model.ButtonInfo{
		Text:    btn02txt,
		Data:    "schedule_status:disable",
		BtnType: model.BtnTypeData,
	}

	btn10 := model.ButtonInfo{
		Text:    "删除上一条消息",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn11 := model.ButtonInfo{
		Text:    btn11txt,
		Data:    "schedule_delete_prev:enable",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "schedule_delete_prev:disable",
		BtnType: model.BtnTypeData,
	}
	btn20 := model.ButtonInfo{
		Text:    "置顶",
		Data:    ":1",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    btn21txt,
		Data:    "schedule_pin:enable",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    btn22txt,
		Data:    "schedule_pin:disable",
		BtnType: model.BtnTypeData,
	}

	btn30 := model.ButtonInfo{
		Text:    "文本内容",
		Data:    "schedule_content",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "媒体图片",
		Data:    "schedule_time",
		BtnType: model.BtnTypeData,
	}
	btn32 := model.ButtonInfo{
		Text:    "按钮链接",
		Data:    "schedule_time",
		BtnType: model.BtnTypeData,
	}

	btn40 := model.ButtonInfo{
		Text:    "重复频率",
		Data:    "schedule_repeat",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "发送时间段",
		Data:    "schedule_time",
		BtnType: model.BtnTypeData,
	}
	btn50 := model.ButtonInfo{
		Text:    "开始日期",
		Data:    "schedule_date_start",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "结束日期",
		Data:    "schedule_date_end",
		BtnType: model.BtnTypeData,
	}
	btnBack := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}

	row0 := []model.ButtonInfo{btn00, btn01, btn02}
	row1 := []model.ButtonInfo{btn10, btn11, btn12}
	row2 := []model.ButtonInfo{btn20, btn21, btn22}
	row3 := []model.ButtonInfo{btn30, btn31, btn32}
	row4 := []model.ButtonInfo{btn40, btn41}
	row5 := []model.ButtonInfo{btn50, btn51}
	backRow := []model.ButtonInfo{btnBack}

	rows := [][]model.ButtonInfo{row0, row1, row2, row3, row4, row5, backRow}
	keyboard := utils.MakeKeyboard(rows)
	utils.ScheduleMsgMenuMarkup = keyboard

	content := updateScheduleMsg()
	utils.SendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)

}

func scheduleStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	msg1 := "✅启用"
	msg2 := "关闭"
	scheduleMsg.Enable = true
	if params == "disable" {
		scheduleMsg.Enable = false
		msg1 = "启用"
		msg2 = "✅关闭"
	}
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[0][1].Text = msg1
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[0][2].Text = msg2
	content := updateScheduleMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ScheduleMsgMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

func scheduleDeletePrevHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	msg1 := "✅是"
	msg2 := "否"
	scheduleMsg.DeletePrevMsg = true
	if params == "disable" {
		scheduleMsg.DeletePrevMsg = false
		msg1 = "是"
		msg2 = "✅否"
	}
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[1][1].Text = msg1
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[1][2].Text = msg2
	content := updateScheduleMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ScheduleMsgMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

func schedulePinHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	msg1 := "✅是"
	msg2 := "否"
	scheduleMsg.Pin = true
	if params == "disable" {
		scheduleMsg.Pin = false
		msg1 = "启用"
		msg2 = "✅关闭"
	}
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[2][1].Text = msg1
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[2][2].Text = msg2
	content := updateScheduleMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ScheduleMsgMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

// 发送频率
func scheduleRepeatMenuHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	hours := []string{"1", "2", "3", "4", "6", "8", "12", "24"}
	minutes := []string{"10", "15", "20", "30"}

	row := []model.ButtonInfo{}
	rows := [][]model.ButtonInfo{}
	tip1 := model.ButtonInfo{
		Text:    "【按小时】",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	tip1Row := []model.ButtonInfo{tip1}
	rows = append(rows, tip1Row)
	//按小时
	for i := 0; i < len(hours); i++ {
		btn := model.ButtonInfo{
			Text:    hours[i],
			Data:    "schedule_repeat_hour:" + hours[i],
			BtnType: model.BtnTypeData,
		}
		row = append(row, btn)
		//达到6个就换行
		if i%4 == 0 && i != 0 {
			rows = append(rows, row)
			row = []model.ButtonInfo{}
		} else if i == 7 {
			rows = append(rows, row)
		}
	}

	tip2 := model.ButtonInfo{
		Text:    "【按分钟】",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	tip2Row := []model.ButtonInfo{tip2}
	rows = append(rows, tip2Row)
	minuteRow := []model.ButtonInfo{}
	//按分钟
	for i := 0; i < len(minutes); i++ {
		btn := model.ButtonInfo{
			Text:    minutes[i],
			Data:    "schedule_repeat_minute:" + minutes[i],
			BtnType: model.BtnTypeData,
		}
		minuteRow = append(minuteRow, btn)
	}
	rows = append(rows, minuteRow)
	btn := model.ButtonInfo{
		Text:    "返回",
		Data:    "schedule_and",
		BtnType: model.BtnTypeData,
	}

	row = []model.ButtonInfo{btn}
	rows = append(rows, row)
	keyboard := utils.MakeKeyboard(rows)
	content := "🕖 定时消息\n👉🏻 选择该消息多久重复一次："
	sendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)
}

// 发送时间段
func scheduleTimeStartMenuHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	//创建24个ButtonInfo，每行6个
	row := []model.ButtonInfo{}
	rows := [][]model.ButtonInfo{}
	for i := 0; i < 24; i++ {
		btn := model.ButtonInfo{
			Text:    strconv.Itoa(i),
			Data:    "schedule_time_start:" + strconv.Itoa(i),
			BtnType: model.BtnTypeData,
		}
		row = append(row, btn)
		//达到6个就换行
		if i%4 == 0 && i != 0 {
			rows = append(rows, row)
			row = []model.ButtonInfo{}
		} else if i == 23 {
			rows = append(rows, row)
		}
	}
	deleteBtn := model.ButtonInfo{
		Text:    "删除已设置的时间段",
		Data:    "schedule_time_delete",
		BtnType: model.BtnTypeData,
	}
	btn := model.ButtonInfo{
		Text:    "返回",
		Data:    "schedule_and",
		BtnType: model.BtnTypeData,
	}

	row = []model.ButtonInfo{deleteBtn, btn}
	rows = append(rows, row)
	keyboard := utils.MakeKeyboard(rows)
	sendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "🕖 定时消息\n设置一个时段，仅在这个时段内发送，下面选项是一天中的24小时，请选择开始时间：", keyboard, bot)
}

func scheduleTimeEndMenuHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	//创建24个ButtonInfo，每行6个
	row := []model.ButtonInfo{}
	rows := [][]model.ButtonInfo{}
	for i := 0; i < 24; i++ {
		btn := model.ButtonInfo{
			Text:    strconv.Itoa(i),
			Data:    "schedule_time_end:" + strconv.Itoa(i),
			BtnType: model.BtnTypeData,
		}
		row = append(row, btn)
		//达到6个就换行
		if i%4 == 0 && i != 0 {
			rows = append(rows, row)
			row = []model.ButtonInfo{}
		} else if i == 23 {
			rows = append(rows, row)
		}
	}
	btn := model.ButtonInfo{
		Text:    "返回",
		Data:    "schedule_and",
		BtnType: model.BtnTypeData,
	}

	row = []model.ButtonInfo{btn}
	rows = append(rows, row)
	keyboard := utils.MakeKeyboard(rows)
	content := fmt.Sprintf("🕖 定时消息\n设置一个时段，仅在这个时段内发送，下面选项是一天中的24小时\n\n已选择开始时间：%d:00\n请选择结束时间：", scheduleMsg.StartHour)
	sendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)
}

func scheduleTimeChangeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	time, _ := strconv.Atoi(params)
	if cmd == "schedule_time_start" {
		scheduleMsg.StartHour = time
		updateScheduleMsg()
		scheduleTimeEndMenuHandler(update, bot)
	} else if cmd == "schedule_time_end" {
		scheduleMsg.EndHour = time
		updateScheduleMsg()
		scheduleMsgMenu(update, bot)
	}
}

// 发送内容
func scheduleContentHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("🕖 定时消息\n当前已设置的文本内容（点击复制）：\n %s \n\n👉🏻 输入你想要设置内容：", scheduleMsg.Text)
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

// 生效日期
func scheduleDateStartHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("🕖 定时消息\n在开启状态下，到达设定时间才会发送消息，请回复开始时间：\n格式：年/月/日 时:分\n例如：%s", utils.CurrentTime())
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	keybord := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("移动已经设置的时间段"),
			tgbotapi.NewKeyboardButton("返回"),
		))

	msg.ReplyMarkup = keybord
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}
	bot.Send(msg)
}

func scheduleDateEndHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("🕖 定时消息\n到达设定时间后自动停止，请回复终止时间：\n格式：年/月/日 时:分\n例如：%s", utils.CurrentTime())
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

// 回复内容处理
func ScheduleAndContentResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	scheduleMsg.Text = update.Message.Text
	content := "添加完成"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "schedule_and",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateScheduleMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func ScheduleDateStartResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	scheduleMsg.StartDate = update.Message.Text
	content := "添加完成"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "schedule_and",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateScheduleMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func ScheduleDateEndResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	scheduleMsg.EndDate = update.Message.Text
	content := "添加完成"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "schedule_and",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateScheduleMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateScheduleMsg() string {
	content := "🕖 定时消息\n\n"
	status_msg := "状态：❌关闭\n"
	if scheduleMsg.Enable {
		status_msg = "状态：✅启用\n"
	}
	repeatMsg := "频率：每60分钟发送一次\n"

	houreMsg := fmt.Sprintf("发送时间段：%d:00 - %d:00\n", scheduleMsg.StartHour, scheduleMsg.EndHour)

	dateMsg := fmt.Sprintf("生效日期：%s - %s\n", scheduleMsg.StartDate, scheduleMsg.EndDate)

	deleteMsg := "删除上一条消息：❌\n"
	if scheduleMsg.DeletePrevMsg {
		deleteMsg = "删除上一条消息：✅\n"
	}

	pinMsg := "消息置顶：❌\n"
	if scheduleMsg.Pin {
		pinMsg = "消息置顶：✅\n"
	}

	meadiaMsg := "媒体图片：❌\n"
	if len(scheduleMsg.Media) > 0 {
		meadiaMsg = "媒体图片：✅\n"
	}
	linkMsg := "按钮链接：❌\n"
	if len(scheduleMsg.Link) > 0 {
		linkMsg = "按钮链接：✅\n"
	}

	text := "文本内容：\n"
	if len(scheduleMsg.Text) > 0 {
		text = "文本内容：" + scheduleMsg.Text + "\n"
	}
	services.SaveModel(&scheduleMsg, scheduleMsg.ChatId)
	content += status_msg + repeatMsg + houreMsg + dateMsg + deleteMsg + pinMsg + meadiaMsg + linkMsg + text
	return content
}

func updateScheduleList() string {
	content := "🕖 定时消息\n设置在群组中每隔几分钟/小时重复发送的消息。\n\n"

	//scheduleMsgs := []model.ScheduleMsg{scheduleMsg}
	status_msg := "❌关闭"
	if scheduleMsg.Enable {
		status_msg = "✅启用"
	}
	repeatMsg := "每60分钟发送一次"

	deleteMsg := "❌"
	if scheduleMsg.DeletePrevMsg {
		deleteMsg = "删除上一条消息：✅"
	}

	pinMsg := "❌"
	if scheduleMsg.Pin {
		pinMsg = "✅"
	}

	text := "文本内容：\n"
	if len(scheduleMsg.Text) > 0 {
		text = "文本内容：" + scheduleMsg.Text + "\n"
	}
	msg := fmt.Sprintf("消息1\n├状态：%s\n├频率：%s\n├发送时间段：%d:00 - %d:00\n├生效日期：%s - %s\n├删除上一条消息：%s\n├消息置顶：%s\n└文本内容：%s\n\n",
		status_msg, repeatMsg, scheduleMsg.StartHour, scheduleMsg.EndHour, scheduleMsg.StartDate, scheduleMsg.EndDate, deleteMsg, pinMsg, text)

	return content + msg
}
