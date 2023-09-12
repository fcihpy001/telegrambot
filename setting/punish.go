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
	"time"
)

var punishment = model.Punishment{}
var class string
var warningSelection = model.SelectInfo{
	Row:    0,
	Column: 0,
	Text:   "1",
}
var actionSelection = model.SelectInfo{
	Row:    5,
	Column: 0,
	Text:   "禁言",
}

func PunishSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	fmt.Println("cmd", cmd)
	fmt.Println("params", params)
	if cmd == "punish_setting_class" {
		class = params
		punishMenu(update, bot)

	} else if cmd == "punish_setting_type" {
		punishTypeHandler(update, bot, params)

	} else if cmd == "punish_setting_count" {
		count, _ := strconv.Atoi(params)
		warningCountHandler(update, bot, count)

	} else if cmd == "punish_setting_action" {
		warningActionHandler(update, bot)

	} else if cmd == "punish_setting_time" {
		muteTimeMenu(update, bot)
	}
}

func punishMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	where := ""
	if class == "flood" {
		where = fmt.Sprintf("flood_setting_id = %d", floodSetting.ID)
	} else if class == "spam" {
		where = fmt.Sprintf("spam_setting_id = %d", spamsSetting.ID)
	} else if class == "prohibited" {
		where = fmt.Sprintf("prohibited_setting_id = %d", prohibitedSetting.ID)
	}
	err := services.GetModelWhere(where, &punishment)

	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/punish.json", &btns)

	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			//返回键盘选项特殊处理
			btn := btnArray[j]
			if btn.Text == "返回" {
				//返回键盘选项
				btn.Data = getBackActionMsg()
			} else {
				btn.Data = btn.Data + ":" + strconv.Itoa(i) + "&" + strconv.Itoa(j)
			}
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.PunishMenuMarkup = keyboard

	//禁言键盘  类型+时长
	rows2 := append(rows[:2], rows[6:]...)
	keyboard2 := utils.MakeKeyboard(rows2)
	utils.PunishMenuMarkup2 = keyboard2

	//仅动作键盘
	rows3 := append(rows[:2], rows[7:]...)
	keyboard3 := utils.MakeKeyboard(rows3)
	utils.PunishMenuMarkup3 = keyboard3

	//要读取用户设置的数据
	content := updatePunishSetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func punishTypeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	switch params {
	case "warn":
		punishment.PunishType = model.PunishTypeWarning
		utils.PunishMenuMarkup.InlineKeyboard[0][0].Text = "✅警告"
		utils.PunishMenuMarkup.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"
		content := updatePunishSetting()
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("statusHandel", err)
		}
	case "mute":
		punishment.PunishType = model.PunishTypeMute
		utils.PunishMenuMarkup2.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup2.InlineKeyboard[0][1].Text = "✅禁言"
		utils.PunishMenuMarkup2.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup2.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup2.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"
		content := updatePunishSetting()
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup2)
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("statusHandel", err)
		}

	case "kick":
		punishment.PunishType = model.PunishTypeKick
		utils.PunishMenuMarkup3.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup3.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup3.InlineKeyboard[0][2].Text = "✅踢出"
		utils.PunishMenuMarkup3.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup3.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"
		content := updatePunishSetting()
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup3)
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("statusHandel", err)
		}

	case "banAndKick":
		punishment.PunishType = model.PunishTypeBanAndKick
		utils.PunishMenuMarkup2.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup2.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup2.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup2.InlineKeyboard[1][0].Text = "✅踢出+封禁"
		utils.PunishMenuMarkup2.InlineKeyboard[1][1].Text = "仅撤回消息+不惩罚"
		content := updatePunishSetting()
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup2)
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("statusHandel", err)
		}

	case "revoke":
		punishment.PunishType = model.PunishTypeRevoke
		utils.PunishMenuMarkup3.InlineKeyboard[0][0].Text = "警告"
		utils.PunishMenuMarkup3.InlineKeyboard[0][1].Text = "禁言"
		utils.PunishMenuMarkup3.InlineKeyboard[0][2].Text = "踢出"
		utils.PunishMenuMarkup3.InlineKeyboard[1][0].Text = "踢出+封禁"
		utils.PunishMenuMarkup3.InlineKeyboard[1][1].Text = "✅仅撤回消息+不惩罚"
		content := updatePunishSetting()
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup3)
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("statusHandel", err)
		}
	}
}

func warningCountHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, count int) {

	//取消以前的选中
	utils.PunishMenuMarkup.InlineKeyboard[3][warningSelection.Column].Text = warningSelection.Text
	//更新选中
	utils.PunishMenuMarkup.InlineKeyboard[3][count-1].Text = "✅" + strconv.Itoa(count)
	//更新选中信息
	warningSelection.Column = count - 1
	warningSelection.Text = strconv.Itoa(count)

	punishment.WarningCount = count
	content := updatePunishSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 达到警告次数后动作
func warningActionHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	//todo 垃圾命名方式，需要修改
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	text := query[1]
	dd := query[2]
	cc := strings.Split(dd, "&")
	col, _ := strconv.Atoi(cc[1])

	//取消以前的选中
	utils.PunishMenuMarkup.InlineKeyboard[actionSelection.Row][actionSelection.Column].Text = actionSelection.Text
	//更新选中
	utils.PunishMenuMarkup.InlineKeyboard[5][col].Text = "✅" + text
	//更新选中信息
	actionSelection.Column = col
	actionSelection.Text = text

	//更新model数据
	if text == "kick" {
		punishment.WarningAfterPunish = model.PunishTypeKick
	} else if text == "banAndKick" {
		punishment.WarningAfterPunish = model.PunishTypeBanAndKick
	} else if text == "mute" {
		punishment.WarningAfterPunish = model.PunishTypeMute
	}

	content := updatePunishSetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.PunishMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updatePunishSetting() string {
	content := "🔇 反垃圾 \n\n惩罚："
	if class == "prohibited" {
		content = "🔇 违禁词 \n\n惩罚："
	} else if class == "flood" {
		content = "🔇 反刷屏 \n\n惩罚："
	} else if class == "spam" {
		content = "📨 反垃圾 \n\n惩罚："
	} else if class == "userCheck" {
		content = "🔦 用户检查 \n\n惩罚："
	}

	//todo 根据class类型分别处理
	actionMsg := "警告 "
	if punishment.PunishType == model.PunishTypeMute {
		actionMsg = "禁言"
	} else if punishment.PunishType == model.PunishTypeKick {
		actionMsg = "踢出"
	} else if punishment.PunishType == model.PunishTypeBanAndKick {
		actionMsg = "踢出+禁言"
	} else if punishment.PunishType == model.PunishTypeRevoke {
		actionMsg = "仅撤回消息+不惩罚"
	} else if punishment.PunishType == model.PunishTypeWarning {
		afterMsg := "禁言"
		if punishment.WarningAfterPunish == model.PunishTypeKick {
			afterMsg = "踢出"
		} else if punishment.WarningAfterPunish == model.PunishTypeBanAndKick {
			afterMsg = "踢出+禁言"
		} else if punishment.WarningAfterPunish == model.PunishTypeMute {
			afterMsg = "禁言"
		}
		actionMsg = fmt.Sprintf("警告%d次后 %s", punishment.WarningCount, afterMsg)
	}

	content = content + actionMsg
	switch class {
	case "spam":
		spamsSetting.WarningCount = punishment.WarningCount
		spamsSetting.Punish = punishment.PunishType
		spamsSetting.MuteTime = punishment.MuteTime
		spamsSetting.BanTime = punishment.BanTime
		spamsSetting.WarningAfterPunish = punishment.WarningAfterPunish
		updateSpamMsg()

	case "flood":
		floodSetting.WarningCount = punishment.WarningCount
		floodSetting.Punish = punishment.PunishType
		floodSetting.MuteTime = punishment.MuteTime
		floodSetting.BanTime = punishment.BanTime
		floodSetting.WarningAfterPunish = punishment.WarningAfterPunish
		updateFloodMsg()

	case "prohibited":
		prohibitedSetting.WarningCount = punishment.WarningCount
		prohibitedSetting.Punish = punishment.PunishType
		prohibitedSetting.MuteTime = punishment.MuteTime
		prohibitedSetting.BanTime = punishment.BanTime
		prohibitedSetting.WarningAfterPunish = punishment.WarningAfterPunish
		updateProhibitedSettingMsg()

	case "userCheck":
		userCheckSetting.WarningCount = punishment.WarningCount
		userCheckSetting.Punish = punishment.PunishType
		userCheckSetting.MuteTime = punishment.MuteTime
		userCheckSetting.BanTime = punishment.BanTime
		userCheckSetting.WarningAfterPunish = punishment.WarningAfterPunish
		updateUserSettingMsg()
	}
	return content
}

// 禁言时长
func muteTimeMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("🔇 违禁词\n\n当前设置：%d分钟 \n👉 输入处罚禁言的时长（分钟，例如：60）：", punishment.MuteTime)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	keybord := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("返回"),
			tgbotapi.NewKeyboardButton("返回2"),
		))

	msg.ReplyMarkup = keybord
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}
	bot.Send(msg)
}

func BanTimeReply(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	time, _ := strconv.Atoi(update.Message.Text)
	punishment.MuteTime = time
	content := "设置成功\n禁言的时长为：" + update.Message.Text + "分钟"
	btn1 := model.ButtonInfo{
		Text:    "️️️⛔️删除已经设置的文本",
		Data:    "group_welcome_text_remove",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "返回",
		Data:    getBackActionMsg(),
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updatePunishSetting()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func getBackActionMsg() string {
	backAction := ""
	if class == "flood" {
		backAction = "flood_setting_menu"
	} else if class == "spam" {
		backAction = "spam_setting_menu"
	} else if class == "prohibited" {
		backAction = "prohibited_setting_menu"
	} else if class == "userCheck" {
		backAction = "user_check_menu"
	}
	return backAction
}

func punishHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, punishment model.Punishment) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	name := update.Message.From.FirstName

	//要返回的结果
	content := ""
	//result := false

	//惩罚记录
	record := model.PunishRecord{}
	record.ChatId = chatId
	record.UserId = userId
	record.Name = name
	record.Reason = punishment.Reason
	record.ReasonType = punishment.ReasonType
	record.WarningCount = 0
	record.MuteTime = 0

	if punishment.PunishType == model.PunishTypeWarning { //警告
		//获取被警告的次数
		where := fmt.Sprintf("chat_id = %d and user_id = %d and reason_type = %d", chatId, userId, punishment.ReasonType)
		_ = services.GetModelWhere(where, &record)
		if record.WarningCount >= punishment.WarningCount { //超出警告次数
			//执行超出警告次数后的逻辑
			if punishment.WarningAfterPunish == model.PunishTypeMute { //禁言
				muteUser(update, bot, punishment.MuteTime*60, userId)
				record.Punish = model.PunishTypeMute
				record.MuteTime = punishment.MuteTime

			} else if punishment.WarningAfterPunish == model.PunishTypeKick { //踢出
				kickUser(update, bot, update.Message.From.ID)
				record.Punish = model.PunishTypeKick

			} else if punishment.WarningAfterPunish == model.PunishTypeBanAndKick { //踢出+封禁
				banUser(update, bot, userId)
				record.Punish = model.PunishTypeBanAndKick
			}
			record.WarningCount = 0
		} else {
			//	发出警告消息
			content = fmt.Sprintf("@%s 您已触发反刷屏规则，警告一次，已被警告%d次", name, record.WarningCount+1)
			if punishment.Reason == "userCheck" {
				content = fmt.Sprintf("@%s 您已触用户规则检查，警告一次，已被警告%d次", name, record.WarningCount+1)
			} else if punishment.Reason == "spam" {
				content = fmt.Sprintf("@%s 您的消息中有不被允许的内容，警告一次，已被警告%d次", name, record.WarningCount+1)
			} else if punishment.Reason == "prohibited" {
				content = fmt.Sprintf("@%s 您所发的消息中含有违禁词，警告一次，已被警告%d次", name, record.WarningCount+1)
			}

			record.WarningCount = record.WarningCount + 1
			record.Punish = model.PunishTypeWarning
		}
		//result = true
	} else if punishment.PunishType == model.PunishTypeMute { //禁言
		muteUser(update, bot, punishment.MuteTime*60, userId)
		record.Punish = model.PunishTypeMute
		record.MuteTime = punishment.MuteTime
		//result = true

	} else if punishment.PunishType == model.PunishTypeKick { //踢出，1天
		kickUser(update, bot, userId)
		record.Punish = model.PunishTypeKick
		//result = true

	} else if punishment.PunishType == model.PunishTypeBan { //封禁，7天
		banUserHandler(update, bot)
		record.Punish = model.PunishTypeMute
		//result = true

	} else if punishment.PunishType == model.PunishTypeRevoke { //撤回
		content = fmt.Sprintf("@%s，系统检测到您存在刷屏行为，请撤回消息", update.Message.From.FirstName)
		record.Punish = model.PunishTypeRevoke
		//result = true

		//return result
	}
	savePunishRecord(bot, chatId, content, &record, punishment.DeleteNotifyMsgTime)
	//return result
}

func savePunishRecord(bot *tgbotapi.BotAPI, chatId int64, content string, record *model.PunishRecord, deleteTime int64) {

	//存储惩罚记录
	services.SaveModel(&record, record.ChatId)
	if len(content) == 0 {
		return
	}

	//对警告类行为，发送提醒消息
	msg := tgbotapi.NewMessage(chatId, content)
	message, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	//需要把这个消息存到记录中，待将来删除
	task := model.Task{
		MessageId:     message.MessageID,
		Type:          "delete",
		OperationTime: time.Now().Add(time.Duration(deleteTime) * time.Minute).Unix(),
	}
	services.SaveModel(&task, chatId)

	mm := tgbotapi.NewDeleteMessage(chatId, message.MessageID)
	bot.Send(mm)
}
