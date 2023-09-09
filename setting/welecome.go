package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var welcomeSetting model.WelcomeSetting

func WelcomeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "welcome_setting_menu" {
		welcomeSettingMenu(update, bot)

	} else if cmd == "welcome_setting_status" {
		welcomeStatusHandler(update, bot, params == "enable")

	} else if cmd == "welcome_setting_delete_prev" {
		welcomeDeletePrevMsgHandler(update, bot, params == "enable")

	} else if cmd == "welcome_setting_type" {
		welcomeTextSettingMenu(update, bot)

	} else if cmd == "welcome_setting_text_remove_menu" {
		welcomeTextDeleteHandler(update, bot)
	}
}

// welcome主菜单
func welcomeSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_ = services.GetModelData(utils.GroupInfo.GroupId, &welcomeSetting)
	welcomeSetting.ChatId = utils.GroupInfo.GroupId

	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/welcome.json", &btns)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			btn := btnArray[j]
			updateWelcomeButtonStatus(&btn)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.GroupWelcomeMarkup = keyboard

	//要读取用户设置的数据
	content := updateWelcomeMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func welcomeStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
	if enable {
		utils.GroupWelcomeMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.GroupWelcomeMarkup.InlineKeyboard[0][2].Text = "关闭"
	} else {
		utils.GroupWelcomeMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.GroupWelcomeMarkup.InlineKeyboard[0][2].Text = "✅关闭"
	}
	welcomeSetting.Enable = enable

	content := updateWelcomeMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.GroupWelcomeMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func welcomeDeletePrevMsgHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, deletePrev bool) {

	if deletePrev {
		utils.GroupWelcomeMarkup.InlineKeyboard[1][1].Text = "✅删除"
		utils.GroupWelcomeMarkup.InlineKeyboard[1][2].Text = "不删"
		welcomeSetting.DeletePrevMsg = true
	} else {
		utils.GroupWelcomeMarkup.InlineKeyboard[1][1].Text = "删除"
		utils.GroupWelcomeMarkup.InlineKeyboard[1][2].Text = "✅不删"
		welcomeSetting.DeletePrevMsg = false
	}
	content := updateWelcomeMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.GroupWelcomeMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 设置欢迎文本
func welcomeTextSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

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
	bot.Send(msg)
}

func WelcomeTextSettingResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	welcomeSetting.WelcomeText = update.Message.Text
	content := "✅设置成功，点击按钮返回"
	btn1 := model.ButtonInfo{
		Text:    "️️️⛔️删除已经设置的文本",
		Data:    "welcome_setting_text_remove",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "返回",
		Data:    "welcome_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateWelcomeMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 删除欢迎文本
func welcomeTextDeleteHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	welcomeSetting.WelcomeText = ""

	content := "✅ 文本内容已删除，点击按钮返回。"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "welcome_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateWelcomeMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateWelcomeMsg() string {
	content := "🎉 进群欢迎\n\n"
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

	content += enableMsg + "\n" + deletePrevMsg + "\n\n自定义欢迎内容：\n" + welcome_media + "\n" + welcome_button + "\n" + welcome_text
	services.SaveModel(&welcomeSetting, utils.GroupInfo.GroupId)
	return content
}

func updateWelcomeButtonStatus(btn *model.ButtonInfo) {
	if btn.Data == "welcome_setting_status" && welcomeSetting.Enable {
		btn.Text = "✅启用"
	} else if btn.Data == "welcome_setting_status" && !welcomeSetting.Enable {
		btn.Text = "✅关闭"
	} else if btn.Data == "welcome_setting_delete_prev" && welcomeSetting.DeletePrevMsg {
		btn.Text = "✅删除"
	} else if btn.Data == "welcome_setting_delete_prev" && !welcomeSetting.DeletePrevMsg {
		btn.Text = "✅不删"
	}
}
