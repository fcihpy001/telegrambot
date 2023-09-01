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

var userCheckSetting model.UserCheck

func UserCheckMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err := services.GetModelData(update.CallbackQuery.Message.Chat.ID, &userCheckSetting)
	fmt.Println("userCheckSetting-query", userCheckSetting)
	userCheckSetting.ChatId = update.CallbackQuery.Message.Chat.ID

	btn11txt := "❌必须设置名字"
	if userCheckSetting.NameCheck {
		btn11txt = "✅必须设置名字"
	}
	btn12txt := "❌必须设置用户名"
	if userCheckSetting.UserNameCheck {
		btn12txt = "✅必须设置用户名"
	}
	btn21txt := "❌必须设置头像"
	if userCheckSetting.IconCheck {
		btn21txt = "✅必须设置置头像"
	}
	btn22txt := "❌必须设置用订阅频道"
	if userCheckSetting.SubScribe {
		btn22txt = "✅必须设置订阅频道"
	}

	btn11 := model.ButtonInfo{
		Text:    btn11txt,
		Data:    "check_name",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "check_username",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    btn21txt,
		Data:    "check_icon",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    btn22txt,
		Data:    "check_channel",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "黑名单列表",
		Data:    "black_user_list",
		BtnType: model.BtnTypeData,
	}
	btn32 := model.ButtonInfo{
		Text:    "添加黑名单",
		Data:    "black_user_add",
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
		Text:    "🏠返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12}
	row2 := []model.ButtonInfo{btn21, btn22}
	row3 := []model.ButtonInfo{btn31, btn32}
	row4 := []model.ButtonInfo{btn41, btn42}
	row5 := []model.ButtonInfo{btn51}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5}
	keyboard := utils.MakeKeyboard(rows)
	utils.UserCheckMenuMarkup = keyboard

	//要读取用户设置的数据
	content := updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func NameCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userCheckSetting.NameCheck = !userCheckSetting.NameCheck
	if userCheckSetting.NameCheck {
		utils.UserCheckMenuMarkup.InlineKeyboard[0][0].Text = "✅必须设置名字"
	} else {
		utils.UserCheckMenuMarkup.InlineKeyboard[0][0].Text = "❌必须设置名字"
	}

	content := updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.UserCheckMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func UserNameCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userCheckSetting.UserNameCheck = !userCheckSetting.UserNameCheck
	if userCheckSetting.UserNameCheck {
		utils.UserCheckMenuMarkup.InlineKeyboard[0][1].Text = "✅必须设置用户名"
	} else {
		utils.UserCheckMenuMarkup.InlineKeyboard[0][1].Text = "❌必须设置用户名"
	}

	content := updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.UserCheckMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func IconCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userCheckSetting.IconCheck = !userCheckSetting.IconCheck
	if userCheckSetting.IconCheck {
		utils.UserCheckMenuMarkup.InlineKeyboard[1][0].Text = "✅必须设置头像"
	} else {
		utils.UserCheckMenuMarkup.InlineKeyboard[1][0].Text = "❌必须设置头像"
	}

	content := updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.UserCheckMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func SubScribeCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userCheckSetting.SubScribe = !userCheckSetting.SubScribe
	if userCheckSetting.SubScribe {
		utils.UserCheckMenuMarkup.InlineKeyboard[1][1].Text = "✅必须订阅频道"
	} else {
		utils.UserCheckMenuMarkup.InlineKeyboard[1][1].Text = "❌必须订阅频道"
	}

	content := updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.UserCheckMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func BlackUserList(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	words := strings.Split(userCheckSetting.BlackUserList, "&")
	fmt.Println("black user", words)
	fmt.Println("black user count", len(words))
	count := len(words)
	if len(words) == 1 && words[0] == "" {
		count = 0
	}
	content := fmt.Sprintf("🔦 用户检查\n\n⛔️ 禁止包含名字   已添加禁止名单：%d条\n\n", count)
	for _, word := range words {
		content = content + fmt.Sprintf("- %s\n", word)
	}

	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_user_check_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func BlackUserAdd(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🔇 黑名单\\n\\n👉请输入要禁止的名字（一行一个）")
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

func BlackUserAddResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if len(userCheckSetting.BlackUserList) > 0 {
		userCheckSetting.BlackUserList = userCheckSetting.BlackUserList + "&" + update.Message.Text
	} else {
		userCheckSetting.BlackUserList = update.Message.Text
	}

	words := strings.Split(userCheckSetting.BlackUserList, "&")

	content := fmt.Sprintf("已添加 %d 个黑名单:\n", len(words))
	for _, word := range words {
		content = fmt.Sprintf("%s\n - %s", content, word)
	}

	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_user_check_setting",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "继续添加",
		Data:    "black_user_add",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)
	updateUserSettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

//func NameContainWordMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
//
//	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
//	_, err := bot.Send(msg)
//	if err != nil {
//		log.Println(err)
//	}
//
//	words := strings.Split(userCheckSetting.NameNotContainWord, "&")
//
//	content := fmt.Sprintf("🔦 用户检查\n\n⛔️ 禁止包含名字   已添加禁止名单：%d条\n\n", len(words))
//	for _, word := range words {
//		content = content + fmt.Sprintf("- %s\n", word)
//	}
//	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
//	keybord := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("返回"),
//		))
//
//	msg.ReplyMarkup = keybord
//	msg.ReplyMarkup = tgbotapi.ForceReply{
//		ForceReply: true,
//	}
//
//	bot.Send(msg)
//}

func NameContainWord(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

}

func updateUserSettingMsg() string {
	content := "🔦 用户检查\n\n在用户进入群组和发送消息时进行检查和屏蔽。\n\n惩罚：警告 3 次后禁言 60 分钟\n\n自动删除提醒消息：10分钟"
	//if replySetting.Enable == false {
	//	content = "💬 关键词回复\n\n当前状态：关闭❌"
	//	return content
	//}
	//fmt.Println("reply_keyworld", replySetting.KeywordReply)
	////enableMsg := "- " + replySetting.KeywordReply[0].KeyWorld
	//
	//enableMsg := "* match world"
	//
	//content = content + enableMsg + "\n" + "\n- 表示精准触发\n * 表示包含触发"

	//services.SaveReplySettings(&replySetting)
	services.SaveModel(&userCheckSetting, userCheckSetting.ChatId)
	return content
}

func UserCheckSetting(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.UserCheckMenuMarkup)
	bot.Send(msg)
}

func GoUserPunishSetting(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.UserCheckMenuMarkup)
	bot.Send(msg)
}
