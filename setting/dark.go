package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var (
	darkModelSetting = model.DarkModelSetting{}
	err              error
)

func darkModelSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err = services.GetModelData(update.CallbackQuery.Message.Chat.ID, &darkModelSetting)
	darkModelSetting.ChatId = update.CallbackQuery.Message.Chat.ID

	btn12txt := "启用"
	btn13txt := "✅关闭"
	if darkModelSetting.Enable {
		btn12txt = "✅启用"
		btn13txt = "关闭"
	}
	btn22txt := "✅全员禁言"
	btn23txt := "禁止媒体"
	if darkModelSetting.BanType == model.BanTypeMedia {
		btn22txt = "全员禁言"
		btn23txt = "✅禁止媒体"
	}
	btn32txt := "通知"
	btn33txt := "✅不通知"
	if darkModelSetting.Notify {
		btn32txt = "✅通知"
		btn33txt = "不通知"
	}

	btn11 := model.ButtonInfo{
		Text:    "状态",
		Data:    "tost",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "dark_model_status:enable",
		BtnType: model.BtnTypeData,
	}
	btn13 := model.ButtonInfo{
		Text:    btn13txt,
		Data:    "dark_model_status:disable",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "模式",
		Data:    "check_icon",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    btn22txt,
		Data:    "dark_model_ban:message",
		BtnType: model.BtnTypeData,
	}
	btn23 := model.ButtonInfo{
		Text:    btn23txt,
		Data:    "dark_model_ban:media",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "通知",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn32 := model.ButtonInfo{
		Text:    btn32txt,
		Data:    "dark_model_notify:enable",
		BtnType: model.BtnTypeData,
	}
	btn33 := model.ButtonInfo{
		Text:    btn33txt,
		Data:    "dark_model_notify:disable",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "设置时间段",
		Data:    "dark_model_time_setting",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12, btn13}
	row2 := []model.ButtonInfo{btn21, btn22, btn23}
	row3 := []model.ButtonInfo{btn31, btn32, btn33}
	row4 := []model.ButtonInfo{btn41}
	row5 := []model.ButtonInfo{btn51}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5}
	keyboard := utils.MakeKeyboard(rows)
	utils.DarkModelMenuMarkup = keyboard
	content := updateDarkSettingMsg()
	sendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)

}

func sendEditMsgMarkup(
	chatID int64,
	messageID int,
	content string,
	replyMarkup tgbotapi.InlineKeyboardMarkup,
	bot *tgbotapi.BotAPI,
) {
	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, content, replyMarkup)
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Println("darkModelSettingMenu", err)
	}
}

func DarkSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "dark_model_setting" {
		darkModelSettingMenu(update, bot)

	} else if cmd == "dark_model_status" {
		statusHandler(update, bot, params)

	} else if cmd == "dark_model_ban" {
		banModelHandler(update, bot, params)

	} else if cmd == "dark_model_notify" {
		notifyHandler(update, bot, params)

	} else if cmd == "dark_model_time_setting" {
		timeStartMenuHandler(update, bot)

	} else if cmd == "dark_model_time2_setting" {
		timeEndMenuHandler(update, bot)

	} else if cmd == "dark_model_time_start" {
		timeSettingStartHandler(update, bot)

	} else if cmd == "dark_model_time_end" {
		timeSettingStartHandler(update, bot)
	}

}

func statusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	statusMsg1 := "✅启用"
	statusMsg2 := "关闭"
	darkModelSetting.Enable = true
	if params == "disable" {
		darkModelSetting.Enable = false
		statusMsg1 = "启用"
		statusMsg2 = "✅关闭"
	}
	utils.DarkModelMenuMarkup.InlineKeyboard[0][1].Text = statusMsg1
	utils.DarkModelMenuMarkup.InlineKeyboard[0][2].Text = statusMsg2
	content := updateDarkSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.DarkModelMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

func banModelHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	statusMsg1 := "✅全员禁言"
	statusMsg2 := "禁止媒体"
	darkModelSetting.BanType = model.BanTypeMessage
	if params == "media" {
		darkModelSetting.BanType = model.BanTypeMedia
		statusMsg1 = "全员禁言"
		statusMsg2 = "✅禁止媒体"
	}
	utils.DarkModelMenuMarkup.InlineKeyboard[1][1].Text = statusMsg1
	utils.DarkModelMenuMarkup.InlineKeyboard[1][2].Text = statusMsg2
	content := updateDarkSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.DarkModelMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

func notifyHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	statusMsg1 := "✅通知"
	statusMsg2 := "不通知"
	darkModelSetting.Notify = true
	if params == "disable" {
		darkModelSetting.Notify = false
		statusMsg1 = "通知"
		statusMsg2 = "✅不通知"
	}
	utils.DarkModelMenuMarkup.InlineKeyboard[2][1].Text = statusMsg1
	utils.DarkModelMenuMarkup.InlineKeyboard[2][2].Text = statusMsg2
	content := updateDarkSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.DarkModelMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("statusHandel", err)
	}
}

func timeStartMenuHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	//创建24个ButtonInfo，每行6个
	row := []model.ButtonInfo{}
	rows := [][]model.ButtonInfo{}
	for i := 0; i < 24; i++ {
		btn := model.ButtonInfo{
			Text:    strconv.Itoa(i),
			Data:    "dark_model_time_start:" + strconv.Itoa(i),
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
		Data:    "dark_model_setting",
		BtnType: model.BtnTypeData,
	}

	row = []model.ButtonInfo{btn}
	rows = append(rows, row)
	keyboard := utils.MakeKeyboard(rows)
	sendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "🌘 夜间模式\n设置每天指定的时段内开始和结束夜间模式，选择开始时间：\n", keyboard, bot)
}

func timeEndMenuHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	//创建24个ButtonInfo，每行6个
	row := []model.ButtonInfo{}
	rows := [][]model.ButtonInfo{}
	for i := 0; i < 24; i++ {
		btn := model.ButtonInfo{
			Text:    strconv.Itoa(i),
			Data:    "dark_model_time_end:" + strconv.Itoa(i),
			BtnType: model.BtnTypeData,
		}
		fmt.Println("buttonInfo", btn)
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
		Data:    "dark_model_setting",
		BtnType: model.BtnTypeData,
	}

	row = []model.ButtonInfo{btn}
	rows = append(rows, row)
	keyboard := utils.MakeKeyboard(rows)
	content := fmt.Sprintf("🌘 夜间模式\n从%d点开始，选择结束时间：\n", darkModelSetting.BanTimeStart)
	sendEditMsgMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard, bot)
}

func timeSettingStartHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	calldata := update.CallbackQuery.Data
	query := strings.Split(calldata, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	time, _ := strconv.Atoi(params)
	if cmd == "dark_model_time_start" {
		darkModelSetting.BanTimeStart = time
		updateDarkSettingMsg()
		timeEndMenuHandler(update, bot)
	} else if cmd == "dark_model_time_end" {
		darkModelSetting.BanTimeEnd = time
		updateDarkSettingMsg()
		darkModelSettingMenu(update, bot)
	}
}

func updateDarkSettingMsg() string {

	content := "🌘 夜间模式\n指定时间自动开启、解除全员禁言。\n\n"
	enableMsg := "状态：✅启用\n"
	if !darkModelSetting.Enable {
		enableMsg = "状态：❌关闭\n"
		content = content + enableMsg
		return content
	}
	banModelMsg := fmt.Sprintf("├从%d - %d 激活 🤫全员禁言\n", darkModelSetting.BanTimeStart, darkModelSetting.BanTimeEnd)
	if darkModelSetting.BanType == model.BanTypeMedia {
		banModelMsg = fmt.Sprintf("├从%d - %d 激活 🤖禁止媒体\n", darkModelSetting.BanTimeStart, darkModelSetting.BanTimeEnd)
	}
	notifyMsg := "└开始和结束通知：✅\n"
	if !darkModelSetting.Notify {
		notifyMsg = "└开始和结束通知：❌\n"
	}
	services.SaveModel(&darkModelSetting, darkModelSetting.ChatId)
	content = content + enableMsg + banModelMsg + notifyMsg
	return content
}
