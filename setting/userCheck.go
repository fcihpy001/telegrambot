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

// 模块入口
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
		subscribeAddMenu(update, bot)

	} else if cmd == "user_check_black_list" {
		blackUserList(update, bot)

	} else if cmd == "user_check_black_add" {
		blackUserAdd(update, bot)

	}
}

// 用户检查菜单
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

// 状态处理-名字检查
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

// 状态处理-用户名检查
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

// 状态处理-头像检查
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

// 状态处理-订阅检查
func subscribeAddMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🔦 用户检查\n\n"+
		"群成员必须订阅指定频道(或加入指定群)后获得发言权限，并且机器人要在该频道(群组)中\n\n"+
		"👉请输入频道或群组地址，格式：https://t.me/[公开链接]")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("返回"),
		))

	msg.ReplyMarkup = keyboard
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
}

func SubscribeAddResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//判断返回的数据是否是以https://t.me开头
	if !strings.HasPrefix(update.Message.Text, "https://t.me/") {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🔦 用户检查\n\n"+
			"格式有误，请重新输入\n\n"+
			"👉请输入频道或群组地址，格式：https://t.me/[公开链接]")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("返回"),
			))

		msg.ReplyMarkup = keyboard
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
		}
		_, err := bot.Send(msg)
		if err != nil {
			return
		}
		return
	}
	//判断当前机器人是否在这个频道中
	content := "✅设置成功"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "user_check_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)
	userCheckSetting.SubScribe = true
	userCheckSetting.ChannelAddr = update.Message.Text
	updateUserSettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 黑名单用户逻辑-列表
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

// 黑名单用户逻辑-添加
func blackUserAdd(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🔇 黑名单\\n\\n👉请输入要禁止的名字(一行一个)")
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

// 黑名单用户逻辑-添加反馈
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

// 配置数据更新
func updateUserSettingMsg() string {
	content := "🔦 用户检查\n\n在用户进入群组和发送消息时进行检查和屏蔽。\n\n"
	punishMsg := "惩罚措施：无\n"
	if len(userCheckSetting.Punish) > 0 {
		if userCheckSetting.Punish == model.PunishTypeWarning {
			punishMsg = fmt.Sprintf("惩罚措施：警告%d次后%s\n", userCheckSetting.WarningCount, utils.PunishActionStr(userCheckSetting.WarningAfterPunish))
		} else {
			punishMsg = fmt.Sprintf("惩罚措施：%s\n", utils.PunishActionStr(userCheckSetting.Punish))
		}
	}
	deleteNotifyMsg := fmt.Sprintf("自动删除提醒消息:%s", utils.TimeStr(userCheckSetting.DeleteNotifyMsgTime))
	content += punishMsg + deleteNotifyMsg
	userCheckSetting.ChatId = utils.GroupInfo.GroupId
	services.SaveModel(&userCheckSetting, utils.GroupInfo.GroupId)
	return content
}

// 菜单按钮初始化显示
func updateUserCheckButtonStatus(btn *model.ButtonInfo) {
	if btn.Data == "user_check_name" && userCheckSetting.NameCheck {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "user_check_username" && userCheckSetting.UserNameCheck {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "user_check_icon" && userCheckSetting.IconCheck {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "user_check_subscribe" && userCheckSetting.SubScribe && len(userCheckSetting.ChannelAddr) > 0 {
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

	content := ""
	//检查用户名
	if setting.UserNameCheck && update.Message.From.UserName == "" {
		content = "没有设置用户名"
	}
	//检查名字
	if setting.NameCheck && update.Message.From.LastName == "" {
		content = "没有设置名字"
	}
	//获取头像信息
	profile, _ := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: update.Message.From.ID,
		Limit:  5,
		Offset: 0,
	})
	if setting.IconCheck && profile.TotalCount < 1 {
		content = "没有设置头像"
	}

	// 检查是否在黑名单中
	if len(setting.BlackUserList) > 0 &&
		len(update.Message.From.UserName) > 0 &&
		strings.Contains(setting.BlackUserList, update.Message.From.FirstName) {
		content = "是黑名单用户"
	}
	if len(content) == 0 {
		return false
	}
	punishment := model.Punishment{
		PunishType:          setting.Punish,
		WarningCount:        setting.WarningCount,
		WarningAfterPunish:  setting.WarningAfterPunish,
		BanTime:             setting.BanTime,
		MuteTime:            setting.MuteTime,
		DeleteNotifyMsgTime: setting.DeleteNotifyMsgTime,
		Reason:              "userCheck",
		ReasonType:          4,
		Content:             content,
	}
	punishHandler(update, bot, punishment)
	return true
}
