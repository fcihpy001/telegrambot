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
var afterSelection = model.SelectInfo{
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
	if class == "flood" {
		punishment.PunishType = floodSetting.Punish
		punishment.WarningCount = floodSetting.WarningCount
		punishment.WarningAfterPunish = floodSetting.WarningAfterPunish
		punishment.BanTime = floodSetting.BanTime
		punishment.DeleteNotifyMsgTime = floodSetting.DeleteNotifyMsgTime
	} else if class == "spam" {
		punishment.PunishType = spamsSetting.Punish
		punishment.WarningCount = spamsSetting.WarningCount
		punishment.WarningAfterPunish = spamsSetting.WarningAfterPunish
		punishment.BanTime = spamsSetting.BanTime
		punishment.DeleteNotifyMsgTime = spamsSetting.DeleteNotifyMsgTime
	} else if class == "prohibited" {
		punishment.PunishType = prohibitedSetting.Punish
		punishment.WarningCount = prohibitedSetting.WarningCount
		punishment.WarningAfterPunish = prohibitedSetting.WarningAfterPunish
		punishment.BanTime = prohibitedSetting.BanTime
		punishment.DeleteNotifyMsgTime = prohibitedSetting.DeleteNotifyMsgTime
	} else if class == "userCheck" {
		punishment.PunishType = userCheckSetting.Punish
		punishment.WarningCount = userCheckSetting.WarningCount
		punishment.WarningAfterPunish = userCheckSetting.WarningAfterPunish
		punishment.BanTime = userCheckSetting.BanTime
		punishment.DeleteNotifyMsgTime = userCheckSetting.DeleteNotifyMsgTime
	}

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
				updatePunishBtn(&btn)
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
	if len(utils.PunishMenuMarkup2.InlineKeyboard) < 1 {
		utils.SendText(update.CallbackQuery.Message.Chat.ID, "请输入/start重新开始", bot)
		return
	}
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
		utils.PunishMenuMarkup2.InlineKeyboard[2][0].Text = "🔇⏱设置禁言时长"
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
		utils.PunishMenuMarkup2.InlineKeyboard[2][0].Text = "🔇⏱设置封禁时长"
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

	if len(utils.PunishMenuMarkup.InlineKeyboard) < 1 {
		return
	}
	//取消以前的选中
	utils.PunishMenuMarkup.InlineKeyboard[3][warningSelection.Column].Text = warningSelection.Text
	//更新选中
	utils.PunishMenuMarkup.InlineKeyboard[3][count-1].Text = "✅" + strconv.Itoa(count)
	//更新选中信息
	warningSelection.Column = count - 1
	warningSelection.Text = strconv.Itoa(count)

	punishment.WarningCount = count
	content := updatePunishSetting()
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
	if len(utils.PunishMenuMarkup.InlineKeyboard) < 1 {
		utils.SendText(update.Message.Chat.ID, "请输入/start重新开始", bot)
		return
	}

	//更新model数据
	if text == "kick" {
		punishment.WarningAfterPunish = model.PunishTypeKick
		afterSelection.Text = "踢出"
	} else if text == "banAndKick" {
		punishment.WarningAfterPunish = model.PunishTypeBanAndKick
		afterSelection.Text = "踢出+封禁"
		utils.PunishMenuMarkup.InlineKeyboard[6][0].Text = "🔇⏱设置封禁时长"
	} else if text == "mute" {
		punishment.WarningAfterPunish = model.PunishTypeMute
		afterSelection.Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[6][0].Text = "🔇⏱设置禁言时长"
	}
	//更新选中
	utils.PunishMenuMarkup.InlineKeyboard[5][col].Text = "✅" + afterSelection.Text
	//更新选中信息
	afterSelection.Column = col

	content := updatePunishSetting()
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		content,
		utils.PunishMenuMarkup)
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
		actionMsg = fmt.Sprintf("警告%d次后 %s", punishment.WarningCount, utils.PunishActionStr(punishment.WarningAfterPunish))
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
	content := ""
	if punishment.PunishType == model.PunishTypeBanAndKick || punishment.WarningAfterPunish == model.PunishTypeBanAndKick {
		content = fmt.Sprintf("🔇 违禁词\n\n当前设置：%d分钟 \n👉 输入处罚封禁的时长（分钟，例如：60）：", punishment.BanTime)
	} else if punishment.PunishType == model.PunishTypeMute || punishment.WarningAfterPunish == model.PunishTypeMute {
		content = fmt.Sprintf("🔇 违禁词\n\n当前设置：%d分钟 \n👉 输入处罚禁言的时长（分钟，例如：60）：", punishment.MuteTime)
	}
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
	content := "设置成功\n禁言的时长为：" + update.Message.Text + "分钟"
	if punishment.PunishType == model.PunishTypeBanAndKick || punishment.WarningAfterPunish == model.PunishTypeBanAndKick {
		punishment.BanTime = time
		content = "设置成功\n封禁的时长为：" + update.Message.Text + "分钟"
	} else if punishment.PunishType == model.PunishTypeMute || punishment.WarningAfterPunish == model.PunishTypeMute {
		punishment.MuteTime = time
		content = "设置成功\n禁言的时长为：" + update.Message.Text + "分钟"
	}
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
	if update.Message.From.IsBot {
		return
	}
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
				MuteUser(chatId, bot, punishment.MuteTime*60, userId)
				record.Punish = model.PunishTypeMute
				record.MuteTime = punishment.MuteTime

			} else if punishment.WarningAfterPunish == model.PunishTypeKick { //踢出
				kickUser(update, bot, update.Message.From.ID)
				record.Punish = model.PunishTypeKick

			} else if punishment.WarningAfterPunish == model.PunishTypeBanAndKick { //踢出+封禁
				banUser(update, bot, userId, uint(punishment.BanTime))
				record.Punish = model.PunishTypeBanAndKick
			}
			record.WarningCount = 0
		} else {
			//	发出警告消息
			content = fmt.Sprintf("@%s 您已触发反刷屏规则:%s，现警告一次，已被警告%d次,警告%d次后会被%s",
				name,
				punishment.Content,
				record.WarningCount+1,
				punishment.WarningCount,
				utils.PunishActionStr(punishment.WarningAfterPunish))
			if punishment.Reason == "userCheck" {
				content = fmt.Sprintf("@%s 您已违法用户检查规则:%s，现警告一次，已被警告%d次,警告%d次后会被%s",
					name,
					punishment.Content,
					record.WarningCount+1,
					punishment.WarningCount,
					utils.PunishActionStr(punishment.WarningAfterPunish))
			} else if punishment.Reason == "spam" {
				content = fmt.Sprintf("@%s 您已违法垃圾消息检查规则:%s，现警告一次，已被警告%d次,警告%d次后会被%s",
					name,
					punishment.Content,
					record.WarningCount+1,
					punishment.WarningCount,
					utils.PunishActionStr(punishment.WarningAfterPunish))
			} else if punishment.Reason == "prohibited" {
				content = fmt.Sprintf("@%s 您所发的消息中含有违禁词，现警告一次，已被警告%d次,警告%d次后会被%s",
					name,
					record.WarningCount+1,
					punishment.WarningCount,
					utils.PunishActionStr(punishment.WarningAfterPunish))
			}
			record.WarningCount = record.WarningCount + 1
			record.Punish = model.PunishTypeWarning
		}

	} else if punishment.PunishType == model.PunishTypeMute { //禁言
		MuteUser(chatId, bot, punishment.MuteTime*60, userId)
		record.Punish = model.PunishTypeMute
		record.MuteTime = punishment.MuteTime

	} else if punishment.PunishType == model.PunishTypeKick { //踢出，1天
		kickUser(update, bot, userId)
		record.Punish = model.PunishTypeKick

	} else if punishment.PunishType == model.PunishTypeBanAndKick { //封禁，7天
		banUser(update, bot, userId, uint(punishment.BanTime))
		record.Punish = model.PunishTypeMute

	} else if punishment.PunishType == model.PunishTypeRevoke { //撤回
		content = fmt.Sprintf("@%s，系统检测到您存在刷屏行为，请撤回消息", update.Message.From.FirstName)
		if punishment.Reason == "userCheck" {
			content = fmt.Sprintf("@%s 您已触用户规则检查,请撤回消息", name)
		} else if punishment.Reason == "spam" {
			content = fmt.Sprintf("@%s 您的消息中有不被允许的内容,请撤回消息", name)
		} else if punishment.Reason == "prohibited" {
			content = fmt.Sprintf("@%s 您所发的消息中含有违禁词,请撤回消息", name)
		}
		record.Punish = model.PunishTypeRevoke
	}
	savePunishRecord(bot, chatId, content, &record, int64(punishment.DeleteNotifyMsgTime))
}

func savePunishRecord(bot *tgbotapi.BotAPI, chatId int64, content string, record *model.PunishRecord, deleteTime int64) {

	//存储惩罚记录
	services.SaveModel(&record, record.ChatId)
	if len(content) == 0 || deleteTime == -1 {
		return
	}

	//对警告类行为，发送提醒消息
	msg := tgbotapi.NewMessage(chatId, content)
	message, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	//需要把这个消息存到记录中，待将来删除
	if deleteTime > 0 {
		task := model.ScheduleDelete{
			ChatId:     chatId,
			MessageId:  message.MessageID,
			DeleteTime: time.Now().Add(time.Duration(deleteTime) * time.Second),
		}
		//保存定时任务
		services.SaveModel(&task, chatId)
	}
}

func updatePunishBtn(btn *model.ButtonInfo) {
	fmt.Println("type:", punishment.PunishType)
	if btn.Data == "punish_setting_type:warn" && punishment.PunishType == model.PunishTypeWarning {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "punish_setting_type:mute" && punishment.PunishType == model.PunishTypeMute {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "punish_setting_type:kick" && punishment.PunishType == model.PunishTypeKick {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "punish_setting_type:banAndKick" && punishment.PunishType == model.PunishTypeBanAndKick {
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "punish_setting_type:revoke" && punishment.PunishType == model.PunishTypeRevoke {
		btn.Text = "✅" + btn.Text
	} else if btn.Text == "1" && punishment.WarningCount == 1 {
		warningSelection.Text = btn.Text
		warningSelection.Column = 0
		btn.Text = "✅" + btn.Text
	} else if btn.Text == "2" && punishment.WarningCount == 2 {
		warningSelection.Text = btn.Text
		warningSelection.Column = 1
		btn.Text = "✅" + btn.Text
	} else if btn.Text == "3" && punishment.WarningCount == 3 {
		warningSelection.Text = btn.Text
		warningSelection.Column = 2
		btn.Text = "✅" + btn.Text
	} else if btn.Text == "4" && punishment.WarningCount == 4 {
		warningSelection.Text = btn.Text
		warningSelection.Column = 3
		btn.Text = "✅" + btn.Text
	} else if btn.Text == "5" && punishment.WarningCount == 5 {
		warningSelection.Text = btn.Text
		warningSelection.Column = 4
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "punish_setting_action:mute" && punishment.WarningAfterPunish == model.PunishTypeMute {
		afterSelection.Text = btn.Text
		afterSelection.Column = 0
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "punish_setting_action:kick" && punishment.WarningAfterPunish == model.PunishTypeKick {
		afterSelection.Text = btn.Text
		afterSelection.Column = 1
		btn.Text = "✅" + btn.Text
	} else if btn.Data == "punish_setting_action:banAndKick" && punishment.WarningAfterPunish == model.PunishTypeBanAndKick {
		afterSelection.Text = btn.Text
		afterSelection.Column = 2
		btn.Text = "✅" + btn.Text
	}
}
