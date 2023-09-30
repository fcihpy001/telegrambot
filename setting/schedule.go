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
	"time"
)

var selectInfo = model.SelectInfo{}
var msgId uint = 0
var msgs []model.ScheduleMsg
var scheduleMessage model.ScheduleMsg

// 处理入口
func ScheduleSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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
		scheduleMessageMenu(update, bot)

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

	} else if cmd == "schedule_delete" {
		ScheduleDelete(update, bot, params)

	} else if cmd == "schedule_modify" {
		scheduleMessageModifyMenu(update, bot, params)

	}
}

// 消息列表菜单
func scheduleMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msgId = 0
	//获取列表
	where := fmt.Sprintf("chat_id = %d", utils.GroupInfo.GroupId)
	messages, _ := services.GetScheduleMsgList(where)
	msgs = messages
	var rows [][]model.ButtonInfo
	for index, _ := range msgs {
		btn1 := model.ButtonInfo{
			Text:    "消息" + strconv.Itoa(index+1),
			Data:    "toast",
			BtnType: model.BtnTypeData,
		}
		btn3 := model.ButtonInfo{
			Text:    "修改",
			Data:    "schedule_modify:" + strconv.Itoa(index),
			BtnType: model.BtnTypeData,
		}
		btn4 := model.ButtonInfo{
			Text:    "删除",
			Data:    "schedule_delete:" + strconv.Itoa(index),
			BtnType: model.BtnTypeData,
		}
		row1 := []model.ButtonInfo{btn1, btn3, btn4}
		rows = append(rows, row1)
	}

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
	//rows = [][]model.ButtonInfo{row1, row2}
	rows = append(rows, row1)
	rows = append(rows, row2)
	keyboard := utils.MakeKeyboard(rows)
	content := updateScheduleList()
	utils.SendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)

}

// 添加时消息菜单
func scheduleMessageMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	where := fmt.Sprintf("chat_id = %d and id = %d", utils.GroupInfo.GroupId, msgId)
	err := services.GetModelWhere(where, &scheduleMessage)
	scheduleMessage.ChatId = utils.GroupInfo.GroupId

	var buttons [][]model.ButtonInfo
	utils.Json2Button2("./config/schedule.json", &buttons)
	fmt.Println(&buttons)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btnArr := buttons[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArr); j++ {
			btn := btnArr[j]
			updateScheduleMessageBtn(&btn)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.ScheduleMsgMenuMarkup = keyboard

	content := updateScheduleMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 修改定时消息菜单
func scheduleMessageModifyMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	index, _ := strconv.Atoi(params)
	message := msgs[index]
	msgId = message.ID
	where := fmt.Sprintf("chat_id = %d and id = %d", utils.GroupInfo.GroupId, msgId)
	err := services.GetModelWhere(where, &scheduleMessage)
	scheduleMessage.ChatId = utils.GroupInfo.GroupId

	var buttons [][]model.ButtonInfo
	utils.Json2Button2("./config/schedule.json", &buttons)
	fmt.Println(&buttons)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btnArr := buttons[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArr); j++ {
			btn := btnArr[j]
			updateScheduleMessageBtn(&btn)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.ScheduleMsgMenuMarkup = keyboard

	content := updateScheduleMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 状态处理
func scheduleStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	msg1 := "✅启用"
	msg2 := "关闭"
	scheduleMessage.Enable = true
	if params == "disable" {
		scheduleMessage.Enable = false
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
	scheduleMessage.DeletePrevMsg = true
	if params == "disable" {
		scheduleMessage.DeletePrevMsg = false
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
	scheduleMessage.Pin = true
	if params == "disable" {
		scheduleMessage.Pin = false
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

// 发送频率菜单及逻辑处理
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
		h, _ := strconv.Atoi(hours[i])

		btn := model.ButtonInfo{
			Text:    hours[i],
			Data:    "schedule_repeat_hour:" + strconv.Itoa(h*60),
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

func repeatValueHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	selectInfo.Text = params

	msg1 := "✅启用"
	msg2 := "关闭"
	scheduleMessage.Enable = true
	if params == "disable" {
		scheduleMessage.Enable = false
		msg1 = "启用"
		msg2 = "✅关闭"
	}
	num, _ := strconv.Atoi(params)
	scheduleMessage.Repeat = num

	utils.ScheduleMsgMenuMarkup.InlineKeyboard[0][1].Text = msg1
	utils.ScheduleMsgMenuMarkup.InlineKeyboard[0][2].Text = msg2
	content := updateScheduleMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ScheduleMsgMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
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
	content := fmt.Sprintf("🕖 定时消息\n设置一个时段，仅在这个时段内发送，下面选项是一天中的24小时\n\n已选择开始时间：%d:00\n请选择结束时间：", scheduleMessage.StartHour)
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
		scheduleMessage.StartHour = time
		updateScheduleMsg()
		scheduleTimeEndMenuHandler(update, bot)
	} else if cmd == "schedule_time_end" {
		scheduleMessage.EndHour = time
		updateScheduleMsg()
		scheduleMessageMenu(update, bot)
	}
}

// 发送内容
func scheduleContentHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("🕖 定时消息\n当前已设置的文本内容（点击复制）：\n %s \n\n👉🏻 输入你想要设置内容：", scheduleMessage.Text)
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
	content := fmt.Sprintf("🕖 定时消息\n在开启状态下，到达设定时间才会发送消息，请回复开始时间：\n格式：年-月-日\n例如：%s", utils.CurrentTime())
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
	content := fmt.Sprintf("🕖 定时消息\n到达设定时间后自动停止，请回复终止时间：\n格式：年-月-日\n例如：%s", utils.CurrentTime())
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
	scheduleMessage.Text = update.Message.Text
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

// 开始日期结果处理
func ScheduleDateStartResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	scheduleMessage.StartDate = update.Message.Text
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

// 结束日期结果处理
func ScheduleDateEndResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	scheduleMessage.EndDate = update.Message.Text
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

// 删除
func ScheduleDelete(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	where := fmt.Sprintf("chat_id = %d", utils.GroupInfo.GroupId)
	messages, _ := services.GetScheduleMsgList(where)
	index, _ := strconv.Atoi(params)

	msg := messages[index]
	wh := fmt.Sprintf("id = %d", msg.ID)
	//删除数据库数据
	_ = services.DeleteModel(&msg, wh)
	//	刷新数据
	scheduleMenu(update, bot)
}

// 更新数据
func updateScheduleMsg() string {
	//msg := "🕖 定时消息\n\n"
	content := fmt.Sprintf("🕖 定时消息\n\n" + scheduleMsgInfo(scheduleMessage))
	if len(scheduleMessage.Text) > 0 && scheduleMessage.ChatId != 0 {
		where := fmt.Sprintf("chat_id = %d and id = '%d'", scheduleMessage.ChatId, msgId)
		scheduleMessage.ExecuteTime = time.Now()
		services.SaveModelWhere(&scheduleMessage, scheduleMessage.ChatId, where)
		msgId = scheduleMessage.ID
	}

	return content
}

func scheduleMsgInfo(msg model.ScheduleMsg) string {
	statusMsg := "❌关闭\n"
	if msg.Enable {
		statusMsg = "✅启用\n"
	}
	repeatMsg := "每60分钟发送一次\n"
	if msg.Repeat > 0 {
		if msg.Repeat < 60 {
			repeatMsg = fmt.Sprintf("每%d分钟发送一次\n", msg.Repeat)
		} else {
			repeatMsg = fmt.Sprintf("每%d小时发送一次\n", msg.Repeat/60)
		}
	}

	deleteMsg := "❌\n"
	if msg.DeletePrevMsg {
		deleteMsg = "✅\n"
	}

	pinMsg := "❌\n"
	if msg.Pin {
		pinMsg = "✅\n"
	}

	text := ""
	if len(msg.Text) > 0 {
		text = msg.Text + "\n"
	}

	content := fmt.Sprintf("消息配置\n"+
		"├状态：%s"+
		"├频率：%s"+
		"├发送时间段：%d:00 - %d:00\n"+
		"├生效日期：%s 至 %s\n"+
		"├删除上一条消息：%s"+
		"├消息置顶：%s"+
		"└文本内容：%s\n\n",
		statusMsg, repeatMsg, msg.StartHour, msg.EndHour, msg.StartDate, msg.EndDate, deleteMsg, pinMsg, text)
	return content
}

func updateScheduleList() string {
	content := "🕖 定时消息\n设置在群组中每隔几分钟/小时重复发送的消息。\n\n"

	where := fmt.Sprintf("chat_id = %d", utils.GroupInfo.GroupId)
	msgs, _ := services.GetScheduleMsgList(where)
	for _, msg := range msgs {
		content += scheduleMsgInfo(msg)
	}
	return content
}

func updateScheduleMessageBtn(btn *model.ButtonInfo) {

	if btn.Data == "schedule_status:enable" && scheduleMessage.Enable {
		btn.Text = "✅启用"
	} else if btn.Data == "schedule_status:disable" && !scheduleMessage.Enable {
		btn.Text = "✅关闭"
	} else if btn.Data == "schedule_delete_prev:enable" && scheduleMessage.DeletePrevMsg {
		btn.Text = "✅是"
	} else if btn.Data == "schedule_delete_prev:disable" && !scheduleMessage.DeletePrevMsg {
		btn.Text = "✅否"
	} else if btn.Data == "schedule_pin:enable" && scheduleMessage.Pin {
		btn.Text = "✅是"
	} else if btn.Data == "schedule_pin:disable" && !scheduleMessage.Pin {
		btn.Text = "✅否"
	}
}

// 定时发消息
func SendMessageTask(bot *tgbotapi.BotAPI) {
	messages, _ := services.GetScheduleMsgList("")
	for _, message := range messages {
		//检查是否达到发送时间
		if !time.Now().After(message.ExecuteTime) {
			continue
		}
		//检查发送开关
		if !message.Enable {
			continue
		}
		//	判断是否在可发送的时间范围内
		if !utils.IsInDateSpan(message.StartDate, message.EndDate) {
			continue
		}
		//删除上次发的消息
		if message.DeletePrevMsg {
			deleteMsg := tgbotapi.NewDeleteMessage(message.ChatId, int(message.MessageId))
			_, _ = bot.Send(deleteMsg)
		}

		//发送消息，并记录messageId
		msg := tgbotapi.NewMessage(message.ChatId, message.Text)
		mm, err := bot.Send(msg)
		if err != nil {
			fmt.Println("err", err)
		}
		message.MessageId = mm.MessageID

		//置顶消息
		if message.Pin {
			PinMessage(message.ChatId, bot, message.MessageId)
		}

		//	更新下次发送时间
		message.ExecuteTime = time.Now().Add(time.Minute * time.Duration(message.Repeat))

		//更新数据库
		services.SaveModel(&message, message.ChatId)
	}
}

// 定时删除消息
func DeleteMessageTask(bot *tgbotapi.BotAPI) {
	//获取所有删除任务
	tasks, _ := services.GetAllDeleteTask()
	//循环tasks
	for _, task := range tasks {
		if time.Now().Before(task.DeleteTime) {
			continue
		}
		msg := tgbotapi.NewDeleteMessage(task.ChatId, task.MessageId)
		mm, err := bot.Send(msg)
		if err != nil {
			fmt.Println("err", err)
		}
		fmt.Println("删除消息成功,更新删除记录", mm.MessageID)
		_ = services.DeleteTask(&task)
	}

}

// 定时检测夜间模式
func CheckDarkTask(bot *tgbotapi.BotAPI) {
	settings, _ := services.GetAllDarkSettings()
	for _, setting := range settings {
		//禁言模式开启且在时间范围内
		if setting.Enable && utils.IsInHoursRange(setting.MuteTimeStart, setting.MuteTimeEnd) {

			//对群组里进行禁言和禁止媒体,计算当前时间到结束时间的时间差
			second := utils.CalculateTimeDifferenceInSeconds(setting.MuteTimeEnd)
			MuteGroup(setting.ChatId, bot, second, setting.MuteType == model.MuteTypeMedia)

			//判断是否发送过
			if setting.OnMessageId != 0 {
				continue
			}

			if !setting.Notify {
				//不通知，且把上一次通知的消息清除掉
				utils.DeleteMessage(setting.ChatId, setting.OffMessageId, bot)
				utils.DeleteMessage(setting.ChatId, setting.OnMessageId, bot)
				continue
			}
			utils.DeleteMessage(setting.ChatId, setting.OffMessageId, bot)
			//组装要发送的消息
			content := fmt.Sprintf("🌘夜间模式开始\n\n❌从现在起禁止发送消息，%d点自动关闭。", setting.MuteTimeEnd)
			if setting.MuteType == model.MuteTypeMedia {
				content = fmt.Sprintf("🌘夜间模式开始\n\n❌从现在起禁止发送媒体消息，%d点自动关闭。", setting.MuteTimeEnd)
			}
			//发送通知消息
			messageId := utils.SendMsg(setting.ChatId, content, bot)
			//更新messageid
			if messageId != 0 {
				setting.OnMessageId = messageId
				setting.OffMessageId = 0
				if setting.ChatId == darkModelSetting.ChatId {
					darkModelSetting.OnMessageId = setting.OnMessageId
					darkModelSetting.OffMessageId = setting.OffMessageId
				}
			}
			//保存配置信息
			services.SaveModel(&setting, setting.ChatId)

		} else {

			MuteGroup(setting.ChatId, bot, 0, setting.MuteType == model.MuteTypeMedia)
			//夜间模式关闭
			if setting.OffMessageId != 0 {
				continue
			}

			if !setting.Notify {
				//不通知，且把上一次通知的消息清除掉
				utils.DeleteMessage(setting.ChatId, setting.OffMessageId, bot)
				utils.DeleteMessage(setting.ChatId, setting.OnMessageId, bot)
				continue
			}
			utils.DeleteMessage(setting.ChatId, setting.OnMessageId, bot)
			services.SaveModel(&setting, setting.ChatId)
			messageId := utils.SendMsg(setting.ChatId, "☀夜间模式关闭，快出来聊天啦。", bot)
			if messageId != 0 {
				setting.OffMessageId = messageId
				setting.OnMessageId = 0
				if setting.ChatId == darkModelSetting.ChatId {
					darkModelSetting.OnMessageId = setting.OnMessageId
					darkModelSetting.OffMessageId = setting.OffMessageId
				}
				services.SaveModel(&setting, setting.ChatId)
			}
		}
	}
}
