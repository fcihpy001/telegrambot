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

var floodSetting model.FloodSetting

func FloodSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

}

func FloodSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err := services.GetModelData(update.CallbackQuery.Message.Chat.ID, &floodSetting)
	floodSetting.ChatId = update.CallbackQuery.Message.Chat.ID
	btn22text := "启用"
	btn23text := "✅关闭"
	if floodSetting.Enable {
		btn22text = "✅启用"
		btn23text = "关闭"
	}

	btn31text := "❌违规后清理消息"
	if floodSetting.DeleteMsg {
		btn31text = "✅违规后清理消息"
	}

	btn11 := model.ButtonInfo{
		Text:    "发送消息条数",
		Data:    "flood_msg_count",
		BtnType: model.BtnTypeData,
	}

	btn12 := model.ButtonInfo{
		Text:    "检查时间间隔",
		Data:    "flood_interval",
		BtnType: model.BtnTypeData,
	}

	btn21 := model.ButtonInfo{
		Text:    "状态",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}

	btn22 := model.ButtonInfo{
		Text:    btn22text,
		Data:    "flood_status_enable",
		BtnType: model.BtnTypeData,
	}

	btn23 := model.ButtonInfo{
		Text:    btn23text,
		Data:    "flood_status_disable",
		BtnType: model.BtnTypeData,
	}

	btn31 := model.ButtonInfo{
		Text:    btn31text,
		Data:    "flood_trigger_delete",
		BtnType: model.BtnTypeData,
	}

	btn41 := model.ButtonInfo{
		Text:    "惩罚设置",
		Data:    "flood_punish_setting",
		BtnType: model.BtnTypeData,
	}

	btn42 := model.ButtonInfo{
		Text:    "自动删除提醒消息",
		Data:    "flood_delete_notify",
		BtnType: model.BtnTypeData,
	}

	btn51 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12}
	row2 := []model.ButtonInfo{btn21, btn22, btn23}
	row3 := []model.ButtonInfo{btn31}
	row4 := []model.ButtonInfo{btn41, btn42}
	row5 := []model.ButtonInfo{btn51}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5}
	keyboard := utils.MakeKeyboard(rows)
	utils.FloodSettingMenuMarkup = keyboard

	content := updateFloodMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func FloodIntervalMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

func FloodMsgCountMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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
	content := "添加完成"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "flood_setting",
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
	content := "添加完成"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "flood_setting",
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

func FloodStatus(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
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

func FloodDeleteMsg(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

	punish_msg := fmt.Sprintf("惩罚：%s %d \n", utils.ActionMap[floodSetting.Punishment.Punish], floodSetting.Punishment.BanTime)

	delete_msg := fmt.Sprintf("自动删除提醒消息：%d分钟", floodSetting.Punishment.DeleteNotifyMsgTime)

	content = content + status_msg + setting_msg + punish_msg + delete_msg
	services.SaveModel(&floodSetting, floodSetting.ChatId)
	return content
}
