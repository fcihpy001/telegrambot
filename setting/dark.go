package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
	"time"
)

var (
	darkModelSetting = model.DarkModelSetting{}
	err              error
)

//todo 禁言时间被误写成ban单词，需要注意

func darkModelSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err = services.GetModelData(utils.GroupInfo.GroupId, &darkModelSetting)
	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/dark.json", &btns)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			btn := btnArray[j]
			updateDarkBtn(&btn)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

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

	darkModelSetting.ChatId = utils.GroupInfo.GroupId
	services.SaveModel(&darkModelSetting, darkModelSetting.ChatId)
	content = content + enableMsg + banModelMsg + notifyMsg
	return content
}

func updateDarkBtn(btn *model.ButtonInfo) {
	if btn.Data == "dark_model_status:enable" && darkModelSetting.Enable {
		btn.Text = "✅启用"
	} else if btn.Data == "dark_model_status:disable" && !darkModelSetting.Enable {
		btn.Text = "✅关闭"
	} else if btn.Data == "dark_model_ban:message" && darkModelSetting.BanType == model.BanTypeMessage {
		btn.Text = "✅全员禁言"
	} else if btn.Data == "dark_model_ban:media" && darkModelSetting.BanType == model.BanTypeMedia {
		btn.Text = "✅禁止媒体"
	} else if btn.Data == "dark_model_notify:enable" && darkModelSetting.Notify {
		btn.Text = "✅通知"
	} else if btn.Data == "dark_model_notify:disable" && !darkModelSetting.Notify {
		btn.Text = "✅不通知"
	}
}

func DarkCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	setting := model.DarkModelSetting{}
	_ = services.GetModelData(chatId, &setting)

	//判断开关
	if !setting.Enable {
		return false
	}
	//	判断时间
	currentHour := time.Now().Hour()
	if currentHour >= setting.BanTimeStart && currentHour < setting.BanTimeEnd {
		// 获取当前时间
		currentTime := time.Now()
		// 定义目标时间
		targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), setting.BanTimeEnd, 0, 0, 0, currentTime.Location())
		secondsDifference := targetTime.Sub(currentTime).Seconds()
		banMember(bot, chatId, int(secondsDifference), userId, setting.BanType == model.BanTypeMedia)
		return true
	}
	return false
}
