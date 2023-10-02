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
			updateFloodBtn(&btn)
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
	if len(utils.FloodSettingMenuMarkup.InlineKeyboard) < 1 {
		utils.SendText(update.CallbackQuery.Message.Chat.ID, "请输入/start重新执行", bot)
		return
	}
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
	if len(utils.FloodSettingMenuMarkup.InlineKeyboard) < 3 {
		utils.SendText(update.CallbackQuery.Message.Chat.ID, "请输入/start重新开始", bot)
		return
	}
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

	actionMsg := "警告"
	if floodSetting.Punish == model.PunishTypeBan {
		actionMsg = "禁言"
	} else if floodSetting.Punish == model.PunishTypeKick {
		actionMsg = "踢出"
	} else if floodSetting.Punish == model.PunishTypeBanAndKick {
		actionMsg = "踢出+禁言"
	} else if floodSetting.Punish == model.PunishTypeRevoke {
		actionMsg = "仅撤回消息+不惩罚"
	} else if floodSetting.Punish == model.PunishTypeWarning {
		actionMsg = fmt.Sprintf("警告%d次后%s", floodSetting.WarningCount, utils.PunishActionStr(floodSetting.WarningAfterPunish))
	}
	punish_msg := fmt.Sprintf("惩罚措施：%s \n", actionMsg)

	delete_msg := fmt.Sprintf("自动删除提醒消息：%d分钟\n", floodSetting.DeleteNotifyMsgTime)

	delete_flood_msg := fmt.Sprint("是否清理违规消息:❌ 否")
	if floodSetting.DeleteMsg {
		delete_flood_msg = fmt.Sprint("是否清理违规消息:✅ 是")
	}
	content = content + status_msg + setting_msg + punish_msg + delete_msg + delete_flood_msg
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

	//确定是否删除违规消息
	if setting.DeleteMsg {

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
		Content:             fmt.Sprintf("在%d秒内发送了至少%d条消息", setting.Interval, setting.MsgCount),
	}
	punishHandler(update, bot, punishment)
	return true
}

func updateFloodBtn(btn *model.ButtonInfo) {
	if btn.Text == "启用" && floodSetting.Enable {
		btn.Text = "✅启用"
	} else if btn.Text == "关闭" && !floodSetting.Enable {
		btn.Text = "✅关闭"
	} else if btn.Text == "违规后清理消息" && floodSetting.DeleteMsg {
		btn.Text = "✅违规后清理消息"
	} else if btn.Text == "违规后清理消息" && !floodSetting.DeleteMsg {
		btn.Text = "❌违规后清理消息"
	}
}
