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

func UserCheckHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]

	if cmd == "user_check_menu" {
		userCheckMenu(update, bot)

	} else if cmd == "user_check_name" {
		nameCheck(update, bot)

	} else if cmd == "user_check_username" {
		userNameCheck(update, bot)

	} else if cmd == "user_check_icon" {
		iconCheck(update, bot)

	} else if cmd == "user_check_subscribe" {
		subScribeCheck(update, bot)

	} else if cmd == "user_check_black_list" {
		blackUserList(update, bot)

	} else if cmd == "user_check_black_add" {
		blackUserAdd(update, bot)

	}
}

func userCheckMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err := services.GetModelData(utils.GroupInfo.GroupId, &userCheckSetting)
	fmt.Println("userCheckSetting-query", userCheckSetting)

	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/userCheck.json", &btns)

	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			btn := btnArray[j]
			updateUserCheckButtonStatus(&btn)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

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

func nameCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

func userNameCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

func iconCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

func subScribeCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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

// 黑名单用户处理
func blackUserList(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	words := strings.Split(userCheckSetting.BlackUserList, "&")

	count := len(words)
	if len(words) == 1 && words[0] == "" {
		count = 0
	}
	content := fmt.Sprintf("🔦 用户检查\n\n⛔️ 禁止包含名字   已添加禁止名单：%d条\n\n", count)
	for _, word := range words {
		content = content + fmt.Sprintf("- %s\n", word)
	}

	btn1 := model.ButtonInfo{
		Text:    "➕添加黑名单",
		Data:    "user_check_black_add",
		BtnType: model.BtnTypeData,
	}

	btn2 := model.ButtonInfo{
		Text:    "返回",
		Data:    "user_check_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateUserSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func blackUserAdd(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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
		Data:    "user_check_menu",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "继续添加",
		Data:    "user_check_black_add",
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
	userCheckSetting.ChatId = utils.GroupInfo.GroupId
	services.SaveModel(&userCheckSetting, utils.GroupInfo.GroupId)
	return content
}

func updateUserCheckButtonStatus(btn *model.ButtonInfo) {
	if btn.Data == "user_check_name" && userCheckSetting.NameCheck {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "user_check_username" && userCheckSetting.UserNameCheck {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "user_check_icon" && userCheckSetting.IconCheck {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "user_check_subscribe" && userCheckSetting.SubScribe {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "user_check_name" && !userCheckSetting.NameCheck {
		btn.Text = "❌" + btn.Text
	} else if btn.Data == "user_check_username" && !userCheckSetting.UserNameCheck {
		btn.Text = "❌" + btn.Text
	} else if btn.Data == "user_check_icon" && !userCheckSetting.IconCheck {
		btn.Text = "❌" + btn.Text
	} else if btn.Data == "user_check_subscribe" && !userCheckSetting.SubScribe {
		btn.Text = "❌" + btn.Text
	}
}

func UserValidateCheck(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	chatId := update.Message.Chat.ID
	setting := model.UserCheck{}
	_ = services.GetModelData(chatId, &setting)

	if setting.UserNameCheck && update.Message.From.UserName == "" {
		content := fmt.Sprintf("@%s 🚫请设置用户名", update.Message.From.FirstName)
		utils.SendText(update.Message.Chat.ID, content, bot)
		return true
	}
	if setting.NameCheck && update.Message.From.LastName == "" {
		content := fmt.Sprintf("@%s 🚫请设置名字", update.Message.From.FirstName)
		utils.SendText(update.Message.Chat.ID, content, bot)
		return true
	}
	//获取头像信息
	profile, _ := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: update.Message.From.ID,
		Limit:  5,
		Offset: 0,
	})
	if setting.IconCheck && profile.TotalCount < 1 {
		content := fmt.Sprintf("🚫@%s 请设置头像", update.Message.From.FirstName)
		utils.SendText(update.Message.Chat.ID, content, bot)
		return true
	}

	// 检查是否在黑名单中
	if len(setting.BlackUserList) > 0 && strings.Contains(setting.BlackUserList, update.Message.From.UserName) {
		content := fmt.Sprintf("🚫@%s 你是黑名单用户，已被禁言", update.Message.From.FirstName)
		utils.SendText(update.Message.Chat.ID, content, bot)
		return true
	}
	return false
}
