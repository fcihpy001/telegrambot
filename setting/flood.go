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

var floodSetting model.FloodSetting

func FloodSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "flood_setting_menu" {
		floodSettingMenu(update, bot)

	} else if cmd == "flood_setting_status" {
		floodStatusHandler(update, bot, params == "enable")

	} else if cmd == "flood_setting_count" {
		floodMsgCountMenu(update, bot)

	} else if cmd == "flood_setting_interval" {
		floodIntervalMenu(update, bot)

	} else if cmd == "flood_setting_delete" {
		floodDeleteMsgHandler(update, bot)

	}
}

func floodSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err := services.GetModelData(utils.GroupInfo.GroupId, &floodSetting)
	floodSetting.ChatId = utils.GroupInfo.GroupId

	var buttons [][]model.ButtonInfo
	utils.Json2Button2("./config/flood.json", &buttons)
	fmt.Println(&buttons)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btnArr := buttons[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArr); j++ {
			btn := btnArr[j]
			if btn.Text == "启用" && floodSetting.Enable {
				btn.Text = "✅启用"
			} else if btn.Text == "关闭" && !floodSetting.Enable {
				btn.Text = "✅关闭"
			}
			if btn.Text == "违规后清理消息" && floodSetting.DeleteMsg {
				btn.Text = "✅违规后清理消息"
			} else if btn.Text == "违规后清理消息" && !floodSetting.DeleteMsg {
				btn.Text = "❌违规后清理消息"
			}
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.FloodSettingMenuMarkup = keyboard

	content := updateFloodMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func floodIntervalMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("当前设置：在 %d秒内发送 %d条消息触发反刷屏\n\n👉 请输入统计发送消息的间隔时间（秒）", floodSetting.Interval, floodSetting.MsgCount)
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

func floodMsgCountMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("当前设置：在 %d秒内发送 %d条消息触发反刷屏\n\n👉 请输入时间内发送消息的最大条数：", floodSetting.Interval, floodSetting.MsgCount)
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

func FloodIntervalResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	count, err := strconv.Atoi(update.Message.Text)

	floodSetting.Interval = count
	content := "✅设置成功"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "flood_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateFloodMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func FloodMsgCountResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	count, err := strconv.Atoi(update.Message.Text)

	floodSetting.MsgCount = count
	content := "✅设置成功"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "flood_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateFloodMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func floodStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
	if enable {
		utils.FloodSettingMenuMarkup.InlineKeyboard[1][1].Text = "✅启用"
		utils.FloodSettingMenuMarkup.InlineKeyboard[1][2].Text = "关闭"
	} else {
		utils.FloodSettingMenuMarkup.InlineKeyboard[1][1].Text = "启用"
		utils.FloodSettingMenuMarkup.InlineKeyboard[1][2].Text = "✅关闭"
	}
	floodSetting.Enable = enable
	content := updateFloodMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.FloodSettingMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func floodDeleteMsgHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	floodSetting.DeleteMsg = !floodSetting.DeleteMsg
	if floodSetting.DeleteMsg {
		utils.FloodSettingMenuMarkup.InlineKeyboard[2][0].Text = "✅违规后清理消息"
	} else {
		utils.FloodSettingMenuMarkup.InlineKeyboard[2][0].Text = "❌违规后清理消息"
	}

	content := updateFloodMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.FloodSettingMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateFloodMsg() string {
	content := "💬 反刷屏\n\n"

	status_msg := "状态：❌ 关闭\n"
	if floodSetting.Enable {
		status_msg = "状态：✅ 开启\n"
	}

	setting_msg := fmt.Sprintf("当前设置：在 %d秒内发送 %d条消息触发反刷屏\n", floodSetting.Interval, floodSetting.MsgCount)

	punish_msg := fmt.Sprintf("惩罚：%s %d \n", utils.ActionMap[floodSetting.Punish], floodSetting.MuteTime)

	delete_msg := fmt.Sprintf("自动删除提醒消息：%d分钟", floodSetting.DeleteNotifyMsgTime)

	content = content + status_msg + setting_msg + punish_msg + delete_msg
	services.SaveModel(&floodSetting, floodSetting.ChatId)
	return content
}

// 刷屏行为检测
func FloodCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	setting := model.FloodSetting{}
	err := services.GetModelData(chatId, &setting)
	if err != nil {
		log.Println(err)
	}
	//统计时间段内,用户发送消息条数排行榜，如果用户发送消息条数超过设置的条数，就触发反刷屏规则
	if !setting.Enable {
		return false
	}

	//检查用户行为
	sendCount := services.MessageCountPeriod(chatId, userId, int64(setting.Interval))
	if sendCount < setting.MsgCount {
		return false
	}
	punishment := model.Punishment{
		PunishType:          setting.Punish,
		WarningCount:        setting.WarningCount,
		WarningAfterPunish:  setting.WarningAfterPunish,
		BanTime:             setting.BanTime,
		MuteTime:            setting.MuteTime,
		DeleteNotifyMsgTime: setting.DeleteNotifyMsgTime,
		Reason:              "flood",
		ReasonType:          3,
		Content:             "",
	}
	punishHandler(update, bot, punishment)
	return true
	//
	////要返回的结果
	//result := false
	//content := ""
	//
	////惩罚记录
	//record := model.PunishRecord{}
	//record.ChatId = chatId
	//record.UserId = userId
	//record.Name = name
	//record.Reason = "flood"
	//record.WarningCount = 0
	//record.MuteTime = 0
	//
	//if setting.Punish == model.PunishTypeWarning { //警告
	//	//	检查警告次数
	//	//获取被警告的次数
	//	where := fmt.Sprintf("chat_id = %d and user_id = %d", chatId, userId)
	//	_ = services.GetModelWhere(where, &record)
	//	if record.WarningCount >= setting.WarningCount { //超出警告次数
	//		//执行超出警告次数后的逻辑
	//		if setting.WarningAfterPunish == model.PunishTypeMute { //禁言
	//			muteUser(update, bot, setting.MuteTime*60, userId)
	//			content = fmt.Sprintf("@%s"+
	//				" 您在 %d 秒内发送了 %d 条消息，"+
	//				"已触发反刷屏规则，将被禁言 %d 分钟",
	//				name, setting.Interval, setting.MsgCount, setting.MuteTime)
	//			record.Punish = model.PunishTypeMute
	//			record.MuteTime = setting.MuteTime
	//
	//		} else if setting.WarningAfterPunish == model.PunishTypeKick { //踢出
	//			kickUser(update, bot, update.Message.From.ID)
	//			record.Punish = model.PunishTypeKick
	//
	//		} else if setting.WarningAfterPunish == model.PunishTypeBanAndKick { //踢出+封禁
	//			banUser(update, bot, userId)
	//			record.Punish = model.PunishTypeBanAndKick
	//		}
	//		record.WarningCount = 0
	//		saveRecord(bot, chatId, content, &record, setting.DeleteNotifyMsgTime)
	//	} else {
	//		//	发出警告消息
	//		content = fmt.Sprintf("@%s 您在 %d 秒内发送了 %d 条消息，已触发反刷屏规则，警告一次，已被警告%d次", name, setting.Interval, setting.MsgCount, record.WarningCount+1)
	//		record.WarningCount = record.WarningCount + 1
	//		record.Punish = model.PunishTypeWarning
	//		saveRecord(bot, chatId, content, &record, setting.DeleteNotifyMsgTime)
	//		return true
	//	}
	//	result = true
	//} else if setting.Punish == model.PunishTypeMute { //禁言
	//	muteUser(update, bot, setting.MuteTime, userId)
	//	record.Punish = model.PunishTypeMute
	//	record.MuteTime = setting.MuteTime
	//	result = true
	//
	//} else if setting.Punish == model.PunishTypeKick { //踢出，1天
	//	kickUser(update, bot, userId)
	//	record.Punish = model.PunishTypeKick
	//	result = true
	//
	//} else if setting.Punish == model.PunishTypeBan { //封禁，7天
	//	banUserHandler(update, bot)
	//	record.Punish = model.PunishTypeMute
	//	result = true
	//
	//} else if setting.Punish == model.PunishTypeRevoke { //撤回
	//	content = fmt.Sprintf("@%s，系统检测到您存在刷屏行为，请撤回消息", update.Message.From.FirstName)
	//	record.Punish = model.PunishTypeRevoke
	//	result = true
	//	saveRecord(bot, chatId, content, &record, setting.DeleteNotifyMsgTime)
	//	return result
	//}
	//return result
}

func saveRecord(bot *tgbotapi.BotAPI, chatId int64, content string, record *model.PunishRecord, deleteTime int64) {

	//存储惩罚记录
	services.SaveModel(&record, record.ChatId)
	if len(content) == 0 {
		return
	}

	//对警告类行为，发送提醒消息
	msg := tgbotapi.NewMessage(chatId, content)
	message, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	//需要把这个消息存到记录中，待将来删除
	task := model.Task{
		MessageId:     message.MessageID,
		Type:          "delete",
		OperationTime: time.Now().Add(time.Duration(deleteTime) * time.Minute),
	}
	services.SaveModel(&task, chatId)
}
