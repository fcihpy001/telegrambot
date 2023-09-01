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

func ProhibitedSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	prohibitedSetting = services.GetProhibitSettings(update.CallbackQuery.Message.Chat.ID)
	prohibitedSetting.ChatId = update.CallbackQuery.Message.Chat.ID
	prohibitedSetting.World = "法轮功&利比亚&台独"
	fmt.Println("prohibite:", prohibitedSetting)
	btn12txt := "启用"
	btn13txt := "✅关闭"
	if prohibitedSetting.Enable {
		btn12txt = "✅启用"
		btn13txt = "关闭"
	}

	btn11 := model.ButtonInfo{
		Text:    "是否启用",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}

	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "prohibitedStatus_enable",
		BtnType: model.BtnTypeData,
	}

	btn13 := model.ButtonInfo{
		Text:    btn13txt,
		Data:    "prohibitedStatus_disable",
		BtnType: model.BtnTypeData,
	}

	btn21 := model.ButtonInfo{
		Text:    "添加违禁词",
		Data:    "prohibited_add_menu",
		BtnType: model.BtnTypeData,
	}

	btn22 := model.ButtonInfo{
		Text:    "删除违禁词",
		Data:    "prohibited_delete_menu",
		BtnType: model.BtnTypeData,
	}

	btn31 := model.ButtonInfo{
		Text:    "列表",
		Data:    "prohibited_list",
		BtnType: model.BtnTypeData,
	}

	btn41 := model.ButtonInfo{
		Text:    "惩罚设置",
		Data:    "prohibited_punish_setting",
		BtnType: model.BtnTypeData,
	}

	btn42 := model.ButtonInfo{
		Text:    "自动删除提醒消息",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn51 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12, btn13}
	row2 := []model.ButtonInfo{btn21, btn22}
	row3 := []model.ButtonInfo{btn31}
	row4 := []model.ButtonInfo{btn41, btn42}
	row5 := []model.ButtonInfo{btn51}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5}

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

func ProhibitedAddMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

func ProhibitedAdd(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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
		Data:    "go_prohibited_setting",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "继续添加",
		Data:    "prohibited_add_menu",
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
		Data:    "go_prohibited_setting",
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

func ProhibitedDeleteMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := "🔇 违禁词\n\n请输入要删除的违禁词（一行一个）："

	btn1 := model.ButtonInfo{
		Text:    "清空违禁词",
		Data:    "prohibited_delete",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_prohibited_setting",
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

func ProhibitedDelete(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	prohibitedSetting.World = ""
	content := "已清空"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_prohibited_setting",
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

func ProhibitedStatus(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {

	if enable {
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][2].Text = "关闭"
		prohibitedSetting.Enable = true
	} else {
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.ProhibiteMenuMarkup.InlineKeyboard[0][2].Text = "✅关闭"
		prohibitedSetting.Enable = false
	}

	content := updateProhibitedSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ProhibiteMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 惩罚设置
func PunishSetting(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	punishType := prohibitedSetting.Punish
	punishMsg1 := "✅警告"
	punishMsg2 := "禁言"
	punishMsg3 := "踢出"
	punishMsg4 := "踢出+封禁"
	punishMsg5 := "仅撤回消息+不惩罚"
	if punishType == model.PunishTypeBan {
		punishMsg1 = "警告"
		punishMsg2 = "✅禁言"
		punishMsg3 = "踢出"
		punishMsg4 = "踢出+封禁"
		punishMsg5 = "仅撤回消息+不惩罚"
	} else if punishType == model.PunishTypeKick {
		punishMsg1 = "警告"
		punishMsg2 = "禁言"
		punishMsg3 = "✅踢出"
		punishMsg4 = "踢出+封禁"
		punishMsg5 = "仅撤回消息+不惩罚"
	} else if punishType == model.PunishTypeBanAndKick {
		punishMsg1 = "警告"
		punishMsg2 = "禁言"
		punishMsg3 = "踢出"
		punishMsg4 = "✅踢出+封禁"
		punishMsg5 = "仅撤回消息+不惩罚"
	} else if punishType == model.PunishTypeRevoke {
		punishMsg1 = "警告"
		punishMsg2 = "禁言"
		punishMsg3 = "踢出"
		punishMsg4 = "踢出+封禁"
		punishMsg5 = "✅仅撤回消息+不惩罚"
	}

	btn11 := model.ButtonInfo{
		Text:    punishMsg1,
		Data:    "prohibit_punish_type1",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    punishMsg2,
		Data:    "prohibit_punish_type2",
		BtnType: model.BtnTypeData,
	}
	btn13 := model.ButtonInfo{
		Text:    punishMsg3,
		Data:    "prohibit_punish_type3",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    punishMsg4,
		Data:    "prohibit_punish_type4",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    punishMsg5,
		Data:    "prohibit_punish_type5",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "警告次数",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "1",
		Data:    "prohibit_warning_count1",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "2",
		Data:    "prohibit_warning_count2",
		BtnType: model.BtnTypeData,
	}
	btn43 := model.ButtonInfo{
		Text:    "3",
		Data:    "prohibit_warning_count3",
		BtnType: model.BtnTypeData,
	}
	btn44 := model.ButtonInfo{
		Text:    "4",
		Data:    "prohibit_warning_count4",
		BtnType: model.BtnTypeData,
	}
	btn45 := model.ButtonInfo{
		Text:    "5",
		Data:    "prohibit_warning_count5",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "达到警告3次后",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn61 := model.ButtonInfo{
		Text:    "禁言",
		Data:    "prohibit_warning_after_action1",
		BtnType: model.BtnTypeData,
	}
	btn62 := model.ButtonInfo{
		Text:    "踢出",
		Data:    "prohibit_warning_after_action2",
		BtnType: model.BtnTypeData,
	}
	btn63 := model.ButtonInfo{
		Text:    "踢出+封禁",
		Data:    "prohibit_warning_after_action3",
		BtnType: model.BtnTypeData,
	}
	btn71 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_prohibited_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12, btn13}
	row2 := []model.ButtonInfo{btn21, btn22}
	row3 := []model.ButtonInfo{btn31}
	row4 := []model.ButtonInfo{btn41, btn42, btn43, btn44, btn45}
	row5 := []model.ButtonInfo{btn51}
	row6 := []model.ButtonInfo{btn61, btn62, btn63}
	row7 := []model.ButtonInfo{btn71}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5, row6, row7}
	keyboard := utils.MakeKeyboard(rows)
	utils.PunishMenuMarkup = keyboard

	//要读取用户设置的数据
	content := "🔇 违禁词\n\n惩罚：警告 4 次后踢出+封禁 60 分钟"
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 惩罚动作
func PunishAction(update *tgbotapi.Update, bot *tgbotapi.BotAPI, punishType model.PunishType) {
	if punishType == model.PunishTypeWarning {
		utils.PunishMenuMarkup.InlineKeyboard[0][0].Text = "✅警告"
		utils.PunishMenuMarkup.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"

	} else if punishType == model.PunishTypeBan {
		utils.PunishMenuMarkup.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup.InlineKeyboard[0][1].Text = "✅禁言"
		utils.PunishMenuMarkup.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"
	} else if punishType == model.PunishTypeKick {
		utils.PunishMenuMarkup.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[0][2].Text = "✅踢出"
		utils.PunishMenuMarkup.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"
	} else if punishType == model.PunishTypeBanAndKick {
		utils.PunishMenuMarkup.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[1][0].Text = "✅踢出+封禁"
		utils.PunishMenuMarkup.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"
	} else if punishType == model.PunishTypeRevoke {
		utils.PunishMenuMarkup.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup.InlineKeyboard[1][1].Text = "✅仅撤回消息+不惩罚"
	}
	prohibitedSetting.Punish = punishType
	content := updatePunishSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 警告次数
func WarningCount(update *tgbotapi.Update, bot *tgbotapi.BotAPI, count int) {
	if count == 1 {
		utils.PunishMenuMarkup.InlineKeyboard[3][0].Text = "✅1"
		utils.PunishMenuMarkup.InlineKeyboard[3][1].Text = "2"
		utils.PunishMenuMarkup.InlineKeyboard[3][2].Text = "3"
		utils.PunishMenuMarkup.InlineKeyboard[3][3].Text = "4"
		utils.PunishMenuMarkup.InlineKeyboard[3][4].Text = "5"

	} else if count == 2 {
		utils.PunishMenuMarkup.InlineKeyboard[3][0].Text = "1"
		utils.PunishMenuMarkup.InlineKeyboard[3][1].Text = "✅2"
		utils.PunishMenuMarkup.InlineKeyboard[3][2].Text = "3"
		utils.PunishMenuMarkup.InlineKeyboard[3][3].Text = "4"
		utils.PunishMenuMarkup.InlineKeyboard[3][4].Text = "5"

	} else if count == 3 {
		utils.PunishMenuMarkup.InlineKeyboard[3][0].Text = "1"
		utils.PunishMenuMarkup.InlineKeyboard[3][1].Text = "2"
		utils.PunishMenuMarkup.InlineKeyboard[3][2].Text = "✅3"
		utils.PunishMenuMarkup.InlineKeyboard[3][3].Text = "4"
		utils.PunishMenuMarkup.InlineKeyboard[3][4].Text = "5"

	} else if count == 4 {
		utils.PunishMenuMarkup.InlineKeyboard[3][0].Text = "1"
		utils.PunishMenuMarkup.InlineKeyboard[3][1].Text = "2"
		utils.PunishMenuMarkup.InlineKeyboard[3][2].Text = "3"
		utils.PunishMenuMarkup.InlineKeyboard[3][3].Text = "✅4"
		utils.PunishMenuMarkup.InlineKeyboard[3][4].Text = "5"

	} else if count == 5 {
		utils.PunishMenuMarkup.InlineKeyboard[3][0].Text = "1"
		utils.PunishMenuMarkup.InlineKeyboard[3][1].Text = "2"
		utils.PunishMenuMarkup.InlineKeyboard[3][2].Text = "3"
		utils.PunishMenuMarkup.InlineKeyboard[3][3].Text = "4"
		utils.PunishMenuMarkup.InlineKeyboard[3][4].Text = "✅5"
	}
	prohibitedSetting.WarningCount = count
	content := updatePunishSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 达到警告次数后动作
func WarningAction(update *tgbotapi.Update, bot *tgbotapi.BotAPI, punishType model.PunishType) {
	if punishType == model.PunishTypeBan {
		utils.PunishMenuMarkup.InlineKeyboard[5][0].Text = "✅禁言"
		utils.PunishMenuMarkup.InlineKeyboard[5][1].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[5][2].Text = "踢出+封禁"
	} else if punishType == model.PunishTypeKick {
		utils.PunishMenuMarkup.InlineKeyboard[5][0].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[5][1].Text = "✅踢出"
		utils.PunishMenuMarkup.InlineKeyboard[5][2].Text = "踢出+封禁"
	} else if punishType == model.PunishTypeBanAndKick {
		utils.PunishMenuMarkup.InlineKeyboard[5][0].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[5][1].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[5][2].Text = "✅踢出+封禁"
	}
	prohibitedSetting.WarningAfterPunish = punishType
	content := updatePunishSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 惩罚时间
func PunishTime(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	punishTime1 := "✅10秒"
	punishTime2 := "60秒"
	punishTime3 := "5分钟"
	punishTime4 := "30分钟"
	punishTime5 := "不删除"
	punishTime6 := "不提醒"
	if prohibitedSetting.BanTime == 60 {
		punishTime1 = "10秒"
		punishTime2 = "✅60秒"
		punishTime3 = "5分钟"
		punishTime4 = "30分钟"
		punishTime5 = "不删除"
		punishTime6 = "不提醒"
	} else if prohibitedSetting.BanTime == 300 {
		punishTime1 = "10秒"
		punishTime2 = "60秒"
		punishTime3 = "✅5分钟"
		punishTime4 = "30分钟"
		punishTime5 = "不删除"
		punishTime6 = "不提醒"
	} else if prohibitedSetting.BanTime == 1800 {
		punishTime1 = "10秒"
		punishTime2 = "60秒"
		punishTime3 = "5分钟"
		punishTime4 = "✅30分钟"
		punishTime5 = "不删除"
		punishTime6 = "不提醒"
	} else if prohibitedSetting.BanTime == 0 {
		punishTime1 = "10秒"
		punishTime2 = "60秒"
		punishTime3 = "5分钟"
		punishTime4 = "30分钟"
		punishTime5 = "✅不删除"
		punishTime6 = "不提醒"
	} else if prohibitedSetting.BanTime == -1 {
		punishTime1 = "10秒"
		punishTime2 = "60秒"
		punishTime3 = "5分钟"
		punishTime4 = "30分钟"
		punishTime5 = "不删除"
		punishTime6 = "✅不提醒"
	}

	btn11 := model.ButtonInfo{
		Text:    punishTime1,
		Data:    "prohibited_ban_time_type1",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    punishTime2,
		Data:    "prohibited_ban_time_type2",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    punishTime3,
		Data:    "prohibited_ban_time_type3",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    punishTime4,
		Data:    "prohibited_ban_time_type4",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    punishTime5,
		Data:    "prohibited_ban_time_type5",
		BtnType: model.BtnTypeData,
	}
	btn32 := model.ButtonInfo{
		Text:    punishTime6,
		Data:    "prohibited_ban_time_type6",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_prohibited_setting",
		BtnType: model.BtnTypeData,
	}

	row1 := []model.ButtonInfo{btn11, btn12}
	row2 := []model.ButtonInfo{btn21, btn22}
	row3 := []model.ButtonInfo{btn31, btn32}
	row4 := []model.ButtonInfo{btn41}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4}
	keyboard := utils.MakeKeyboard(rows)
	utils.PunishTimeMarkup = keyboard

	//要读取用户设置的数据
	content := "🔇 违禁词\n\n群成员触发🔇 违禁词时，机器人发出的提醒消息在多少时间后自动删除"
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishTimeMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 惩罚时间
func PunishTimeType(update *tgbotapi.Update, bot *tgbotapi.BotAPI, ban_time_type model.BanTimeType) {

	if ban_time_type == model.BanTimeType1 {
		utils.PunishTimeMarkup.InlineKeyboard[0][0].Text = "✅10秒"
		utils.PunishTimeMarkup.InlineKeyboard[0][1].Text = "60秒"
		utils.PunishTimeMarkup.InlineKeyboard[1][0].Text = "5分钟"
		utils.PunishTimeMarkup.InlineKeyboard[1][1].Text = "30分钟"
		utils.PunishTimeMarkup.InlineKeyboard[2][0].Text = "不删除"
		utils.PunishTimeMarkup.InlineKeyboard[2][1].Text = "不提醒"
		prohibitedSetting.DeleteNotifyMsgTime = 10
	} else if ban_time_type == model.BanTimeType2 {
		utils.PunishTimeMarkup.InlineKeyboard[0][0].Text = "10秒"
		utils.PunishTimeMarkup.InlineKeyboard[0][1].Text = "✅60秒"
		utils.PunishTimeMarkup.InlineKeyboard[1][0].Text = "5分钟"
		utils.PunishTimeMarkup.InlineKeyboard[1][1].Text = "30分钟"
		utils.PunishTimeMarkup.InlineKeyboard[2][0].Text = "不删除"
		utils.PunishTimeMarkup.InlineKeyboard[2][1].Text = "不提醒"
		prohibitedSetting.DeleteNotifyMsgTime = 60
	} else if ban_time_type == model.BanTimeType3 {
		utils.PunishTimeMarkup.InlineKeyboard[0][0].Text = "10秒"
		utils.PunishTimeMarkup.InlineKeyboard[0][1].Text = "60秒"
		utils.PunishTimeMarkup.InlineKeyboard[1][0].Text = "✅5分钟"
		utils.PunishTimeMarkup.InlineKeyboard[1][1].Text = "30分钟"
		utils.PunishTimeMarkup.InlineKeyboard[2][0].Text = "不删除"
		utils.PunishTimeMarkup.InlineKeyboard[2][1].Text = "不提醒"
		prohibitedSetting.DeleteNotifyMsgTime = 300
	} else if ban_time_type == model.BanTimeType4 {
		utils.PunishTimeMarkup.InlineKeyboard[0][0].Text = "10秒"
		utils.PunishTimeMarkup.InlineKeyboard[0][1].Text = "60秒"
		utils.PunishTimeMarkup.InlineKeyboard[1][0].Text = "5分钟"
		utils.PunishTimeMarkup.InlineKeyboard[1][1].Text = "✅30分钟"
		utils.PunishTimeMarkup.InlineKeyboard[2][0].Text = "不删除"
		utils.PunishTimeMarkup.InlineKeyboard[2][1].Text = "不提醒"
		prohibitedSetting.DeleteNotifyMsgTime = 1800
	} else if ban_time_type == model.BanTimeType5 {
		utils.PunishTimeMarkup.InlineKeyboard[0][0].Text = "10秒"
		utils.PunishTimeMarkup.InlineKeyboard[0][1].Text = "60秒"
		utils.PunishTimeMarkup.InlineKeyboard[1][0].Text = "5分钟"
		utils.PunishTimeMarkup.InlineKeyboard[1][1].Text = "30分钟"
		utils.PunishTimeMarkup.InlineKeyboard[2][0].Text = "✅不删除"
		utils.PunishTimeMarkup.InlineKeyboard[2][1].Text = "不提醒"
		prohibitedSetting.DeleteNotifyMsgTime = 0
	} else if ban_time_type == model.BanTimeType6 {
		utils.PunishTimeMarkup.InlineKeyboard[0][0].Text = "10秒"
		utils.PunishTimeMarkup.InlineKeyboard[0][1].Text = "60秒"
		utils.PunishTimeMarkup.InlineKeyboard[1][0].Text = "5分钟"
		utils.PunishTimeMarkup.InlineKeyboard[1][1].Text = "30分钟"
		utils.PunishTimeMarkup.InlineKeyboard[2][0].Text = "不删除"
		utils.PunishTimeMarkup.InlineKeyboard[2][1].Text = "✅不提醒"
		prohibitedSetting.DeleteNotifyMsgTime = -1
	}
	updateProhibitedSettingMsg()
	content := "🔇 违禁词\n\n群成员触发🔇 违禁词时，机器人发出的提醒消息在多少时间后自动删除"
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishTimeMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func GoProhibitedSetting(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := updateProhibitedSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ProhibiteMenuMarkup)
	bot.Send(msg)
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
		actionMsg = fmt.Sprintf("警告%d次后 %s", prohibitedSetting.WarningCount, actionMap[prohibitedSetting.WarningAfterPunish])
	}
	deleteNotifyMsg := "\n自动删除提醒消息：关闭"
	if prohibitedSetting.DeleteNotifyMsgTime > 0 {
		deleteNotifyMsg = fmt.Sprintf("\n自动删除提醒消息：%d ", prohibitedSetting.DeleteNotifyMsgTime)
	} else if prohibitedSetting.DeleteNotifyMsgTime == -1 {
		deleteNotifyMsg = "\n自动删除提醒消息：不提醒"
	} else if prohibitedSetting.DeleteNotifyMsgTime == 0 {
		deleteNotifyMsg = "\n自动删除提醒消息：不删除"
	}

	content = content + enableMsg + actionMsg + deleteNotifyMsg
	services.SaveProhibitSettings(&prohibitedSetting)
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
	notifyTimeMap = map[model.BanTimeType]string{
		model.BanTimeType1: "10秒",
		model.BanTimeType2: "60秒",
		model.BanTimeType3: "5分钟",
		model.BanTimeType4: "30分钟",
		model.BanTimeType5: "不删除",
		model.BanTimeType6: "不提醒",
	}
)

func updatePunishSettingMsg() string {
	content := "🔇 违禁词\n\n惩罚："
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
		actionMsg = fmt.Sprintf("警告%d次后 %s", prohibitedSetting.WarningCount, actionMap[prohibitedSetting.WarningAfterPunish])
	}

	content = content + actionMsg
	services.SaveProhibitSettings(&prohibitedSetting)
	return content
}
