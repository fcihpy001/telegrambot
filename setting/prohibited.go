package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var prohibitedSetting model.ProhibitedSetting

// ProhibitedSettingHandler 违禁词处理逻辑入口
func ProhibitedSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_ = services.GetModelData(utils.GroupInfo.GroupId, &prohibitedSetting)
	scheduleMessage.ChatId = update.CallbackQuery.Message.Chat.ID

	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	fmt.Println(query)
	if cmd == "prohibited_setting_menu" { //违禁词设置主菜单
		prohibitedSettingMenu(update, bot)

	} else if cmd == "prohibited_status" { //违禁词开关
		prohibitedStatus(update, bot, params == "enable")

	} else if cmd == "prohibited_list" { //违禁词列表
		ProhibitedList(update, bot)

	} else if cmd == "prohibited_add" { //违禁词添加
		prohibitedAddMenu(update, bot)

	} else if cmd == "prohibited_delete" { //违禁词删除
		prohibitedDeleteMenu(update, bot)

	} else if cmd == "punish_setting_class" { //违禁词惩罚
		punishMenu(update, bot)

	} else if cmd == "delete_notify_menu" { //违禁词警告
		DeleteNotifyHandler(update, bot)
	}

}

// 违禁词主菜单
func prohibitedSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	_ = services.GetModelData(utils.GroupInfo.GroupId, &prohibitedSetting)
	prohibitedSetting.ChatId = utils.GroupInfo.GroupId

	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/prohibited.json", &btns)

	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			btn := btnArray[j]
			if btn.Data == "prohibited_status:enable" && prohibitedSetting.Enable {
				btn.Text = "✅启用"
			} else if btn.Data == "prohibited_status:disable" && !prohibitedSetting.Enable {
				btn.Text = "✅关闭"
			}
			row = append(row, btn)
		}
		rows = append(rows, row)
	}
	keyboard := utils.MakeKeyboard(rows)
	utils.ProhibiteMenuMarkup = keyboard

	//要读取用户设置的数据
	content := updateProhibitedSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 违禁词添加菜单
func prohibitedAddMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🔇 违禁词\\n\\n👉请输入添加的违禁词（一行一个）")
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

// 违禁词添加结果
func ProhibitedAddResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if len(prohibitedSetting.World) > 0 {
		prohibitedSetting.World = prohibitedSetting.World + "&" + update.Message.Text
	} else {
		prohibitedSetting.World = update.Message.Text
	}

	words := strings.Split(prohibitedSetting.World, "&")

	content := fmt.Sprintf("已添加 %d 个违禁词:\n", len(words))
	for _, word := range words {
		content = fmt.Sprintf("%s\n - %s", content, word)
	}

	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "prohibited_setting_menu",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "继续添加",
		Data:    "prohibited_add",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)
	updateProhibitedSettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	//msg := tgbotapi.NewEditMessageTextAndMarkup(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func ProhibitedList(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	content := "违禁词列表：\n[空]"
	if len(prohibitedSetting.World) > 0 {
		strs := strings.Split(prohibitedSetting.World, "&")
		content = "违禁词列表：\n"
		for _, str := range strs {
			content = fmt.Sprintf("%s\n%s", content, str)
		}
	}
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "prohibited_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 违禁词删除菜单
func prohibitedDeleteMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := "🔇 违禁词\n\n请输入要删除的违禁词（一行一个）："

	btn1 := model.ButtonInfo{
		Text:    "清空违禁词",
		Data:    "prohibited_delete",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "返回",
		Data:    "prohibited_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func ProhibitedDeleteResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	prohibitedSetting.World = ""
	content := "已清空"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "prohibited_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 违禁词状态处理
func prohibitedStatus(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {

	prohibitedSetting.Enable = enable
	if prohibitedSetting.Enable {
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][2].Text = "关闭"
	} else {
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][2].Text = "✅关闭"
	}

	content := updateProhibitedSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ProhibiteMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateProhibitedSettingMsg() string {
	content := "🔇 违禁词\n\n"
	enableMsg := "当前状态：关闭❌\n"
	if prohibitedSetting.Enable {
		enableMsg = "当前状态：启用✅\n"
	}
	actionMsg := "警告"
	if prohibitedSetting.Punish == model.PunishTypeBan {
		actionMsg = "禁言"
	} else if prohibitedSetting.Punish == model.PunishTypeKick {
		actionMsg = "踢出"
	} else if prohibitedSetting.Punish == model.PunishTypeBanAndKick {
		actionMsg = "踢出+禁言"
	} else if prohibitedSetting.Punish == model.PunishTypeRevoke {
		actionMsg = "仅撤回消息+不惩罚"
	} else if prohibitedSetting.Punish == model.PunishTypeWarning {
		actionMsg = fmt.Sprintf("警告%d次后%s", prohibitedSetting.WarningCount, actionMap[prohibitedSetting.WarningAfterPunish])
	}
	deleteNotifyMsg := "\n自动删除提醒消息：关闭"
	if prohibitedSetting.DeleteNotifyMsgTime > 0 {
		deleteNotifyMsg = fmt.Sprintf("\n自动删除提醒消息：%s后", utils.TimeStr(prohibitedSetting.DeleteNotifyMsgTime))
	} else if prohibitedSetting.DeleteNotifyMsgTime == -1 {
		deleteNotifyMsg = "\n自动删除提醒消息：不提醒"
	} else if prohibitedSetting.DeleteNotifyMsgTime == 0 {
		deleteNotifyMsg = "\n自动删除提醒消息：不删除"
	}

	content = content + enableMsg + "惩罚措施:" + actionMsg + deleteNotifyMsg
	services.SaveModel(&prohibitedSetting, prohibitedSetting.ChatId)
	return content
}

var (
	actionMap = map[model.PunishType]string{
		model.PunishTypeWarning:    "警告",
		model.PunishTypeBan:        "禁言",
		model.PunishTypeKick:       "踢出",
		model.PunishTypeBanAndKick: "踢出+封禁",
		model.PunishTypeRevoke:     "仅撤回消息+不惩罚",
	}
)

//func updatePunishSettingMsg() string {
//	content := "🔇 违禁词\n\n惩罚："
//	actionMsg := "警告"
//
//	if prohibitedSetting.Punish == model.PunishTypeBan {
//		actionMsg = "禁言"
//	} else if prohibitedSetting.Punish == model.PunishTypeKick {
//		actionMsg = "踢出"
//	} else if prohibitedSetting.Punish == model.PunishTypeBanAndKick {
//		actionMsg = "踢出+禁言"
//	} else if prohibitedSetting.Punish == model.PunishTypeRevoke {
//		actionMsg = "仅撤回消息+不惩罚"
//	} else if prohibitedSetting.Punish == model.PunishTypeWarning {
//		actionMsg = fmt.Sprintf("惩罚措施:警告%d次后%s", prohibitedSetting.WarningCount, actionMap[prohibitedSetting.WarningAfterPunish])
//	}
//
//	content = content + "惩罚措施:" + actionMsg
//	services.SaveModel(&prohibitedSetting, prohibitedSetting.ChatId)
//	return content
//}

// 过滤违禁词
func ProhibitedCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	messageText := update.Message.Text
	chatId := update.Message.Chat.ID
	//获取数据库中的违禁词列表
	setting := model.ProhibitedSetting{}
	_ = services.GetModelData(chatId, &setting)
	if setting.Enable == false {
		return false
	}

	if !strings.Contains(setting.World, messageText) {
		return false
	}
	punishment := model.Punishment{
		PunishType:          setting.Punish,
		WarningCount:        setting.WarningCount,
		WarningAfterPunish:  setting.WarningAfterPunish,
		BanTime:             setting.BanTime,
		MuteTime:            setting.MuteTime,
		DeleteNotifyMsgTime: setting.DeleteNotifyMsgTime,
		Reason:              "prohibited",
		ReasonType:          1,
		Content:             "",
	}
	punishHandler(update, bot, punishment)
	return true
}
