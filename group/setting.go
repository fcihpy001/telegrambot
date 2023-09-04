package group

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var welcomeSetting model.WelcomeSetting

func (mgr *GroupManager) group_welcome_setting(update *tgbotapi.Update) {
	//从数据库中获取welecome setting
	chatId := update.CallbackQuery.Message.Chat.ID
	log.Println("welcomeSetting:", chatId)
	err := services.GetModelData(utils.GroupInfo.GroupId, &welcomeSetting)
	welcomeSetting.ChatId = utils.GroupInfo.GroupId

	btn12txt := "启用"
	btn13txt := "✅关闭"
	if welcomeSetting.Enable {
		btn12txt = "✅启用"
		btn13txt = "关闭"
	}

	btn22txt := "删除"
	btn23txt := "✅不删"
	if welcomeSetting.DeletePrevMsg {
		btn22txt = "✅删除"
		btn23txt = "不删"
	}

	btn11 := model.ButtonInfo{
		Text:    "是否启用",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "group_welcomeSettingEnable",
		BtnType: model.BtnTypeData,
	}
	btn13 := model.ButtonInfo{
		Text:    btn13txt,
		Data:    "group_welcomeSettingDisable",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "删除上条消息",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    btn22txt,
		Data:    "group_welcome_DeletePrevMsg_enable",
		BtnType: model.BtnTypeData,
	}
	btn23 := model.ButtonInfo{
		Text:    btn23txt,
		Data:    "group_welcome_DeletePrevMsg_disable",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "🦁自定义欢迎内容🦁",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "🦚文本内容",
		Data:    "group_welcome_setting_text",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "🍇媒体图片",
		Data:    "group_welcome_setting_media",
		BtnType: model.BtnTypeData,
	}
	btn43 := model.ButtonInfo{
		Text:    "🍵链接按钮",
		Data:    "group_welcome_setting_button",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "🏠返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12, btn13}
	row2 := []model.ButtonInfo{btn21, btn22, btn23}
	row3 := []model.ButtonInfo{btn31}
	row4 := []model.ButtonInfo{btn41, btn42, btn43}
	row5 := []model.ButtonInfo{btn51}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5}
	keyboard := utils.MakeKeyboard(rows)
	utils.GroupWelcomeMarkup = keyboard

	//要读取用户设置的数据
	content := updateMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *GroupManager) group_welcomeSettingStatus(update *tgbotapi.Update, enable bool) {

	if enable {
		utils.GroupWelcomeMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.GroupWelcomeMarkup.InlineKeyboard[0][2].Text = "关闭"
	} else {
		utils.GroupWelcomeMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.GroupWelcomeMarkup.InlineKeyboard[0][2].Text = "✅关闭"
	}
	welcomeSetting.Enable = enable

	content := updateMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.GroupWelcomeMarkup)
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *GroupManager) welcomeSettingDeletePrevMsg(update *tgbotapi.Update, deletePrev bool) {

	if deletePrev {
		utils.GroupWelcomeMarkup.InlineKeyboard[1][1].Text = "✅删除"
		utils.GroupWelcomeMarkup.InlineKeyboard[1][2].Text = "不删"
		welcomeSetting.DeletePrevMsg = true
	} else {
		utils.GroupWelcomeMarkup.InlineKeyboard[1][1].Text = "删除"
		utils.GroupWelcomeMarkup.InlineKeyboard[1][2].Text = "✅不删"
		welcomeSetting.DeletePrevMsg = false
	}
	content := updateMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.GroupWelcomeMarkup)
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (mgr *GroupManager) welcomeTextSetting(update *tgbotapi.Update) {

	content := "👉 输入要设置的新成员入群欢迎内容，占位符中%s代替，如：👏👏👏 热烈欢迎 %s 加入 %s"
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	keybord := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("返回"),
		))

	msg.ReplyMarkup = keybord
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}
	mgr.bot.Send(msg)

	//content := "👉 输入要设置的新成员入群欢迎内容，占位符中%s代替，如：👏👏👏 热烈欢迎 %s 加入 %s"
	//if len(welcomeSetting.WelcomeText) > 0 {
	//	content = fmt.Sprintf("当前设置的文本(长按下方文字复制)：\n%s\n\n\n👉 输入要设置的新成员入群欢迎内容，占位符中%s代替，如：👏👏👏 热烈欢迎 %s 加入 %s", welcomeSetting.WelcomeText)
	//	rows = [][]model.ButtonInfo{row1, row2}
	//}
	//msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	//keybord := tgbotapi.NewReplyKeyboard(
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton("⛔️删除已经设置的文本"),
	//		tgbotapi.NewKeyboardButton("返回"),
	//	))
	//
	//msg.ReplyMarkup = keybord
	//msg.ReplyMarkup = tgbotapi.ForceReply{
	//	ForceReply: true,
	//}
	//mgr.bot.Send(msg)

	//btn11 := model.ButtonInfo{
	//	Text:    "⛔️删除已经设置的文本",
	//	Data:    "group_welcome_text_remove",
	//	BtnType: model.BtnTypeData,
	//}
	//btn21 := model.ButtonInfo{
	//	Text:    "返回",
	//	Data:    "group_welcome_setting",
	//	BtnType: model.BtnTypeData,
	//}
	//
	//row1 := []model.ButtonInfo{btn11}
	//row2 := []model.ButtonInfo{btn21}
	//rows := [][]model.ButtonInfo{row2}
	//content := "👉 输入要设置的新成员入群欢迎内容，占位符中%s代替，如：👏👏👏 热烈欢迎 %s 加入 %s"
	//if len(welcomeSetting.WelcomeText) > 0 {
	//	content = fmt.Sprintf("当前设置的文本(长按下方文字复制)：\n%s\n\n\n👉 输入要设置的新成员入群欢迎内容，占位符中%s代替，如：👏👏👏 热烈欢迎 %s 加入 %s", welcomeSetting.WelcomeText)
	//	rows = [][]model.ButtonInfo{row1, row2}
	//}
	//keyboard := utils.MakeKeyboard(rows)
	//msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	//
	//_, err := mgr.bot.Send(msg)
	//if err != nil {
	//	log.Println(err)
	//}
}

func WelcomeTextSettingResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	welcomeSetting.WelcomeText = update.Message.Text
	content := "添加完成"
	btn1 := model.ButtonInfo{
		Text:    "️️️⛔️删除已经设置的文本",
		Data:    "group_welcome_text_remove",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "返回",
		Data:    "group_welcome_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateMsg() string {
	content := "🎉 进群欢迎"
	enableMsg := "当前状态：关闭 ❌"
	if welcomeSetting.Enable {
		enableMsg = "当前状态：开启 ✅"
	}
	deletePrevMsg := "删除上条消息：❌"
	if welcomeSetting.DeletePrevMsg {
		deletePrevMsg = "删除上条消息：✅"
	}

	welcome_media := "┌📸 媒体图片:❌"
	welcome_button := "├🔠 链接按钮:❌"
	welcome_text := "└📄 文本内容:❌"
	if len(welcomeSetting.WelcomeText) > 0 {
		welcome_text = "└📄 文本内容: " + welcomeSetting.WelcomeText
	}
	if len(welcomeSetting.WelcomeButton) > 0 {
		welcome_button = "├🔠 链接按钮:✅"
	}
	if len(welcomeSetting.WelcomeMedia) > 0 {
		welcome_media = "📸 媒体图片:✅"
	}

	content = "进群欢迎\n\n" + enableMsg + "\n" + deletePrevMsg + "\n\n自定义欢迎内容：\n" + welcome_media + "\n" + welcome_button + "\n" + welcome_text
	services.SaveModel(&welcomeSetting, utils.GroupInfo.GroupId)
	return content
}
