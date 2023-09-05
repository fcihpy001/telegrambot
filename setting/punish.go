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
		banTimeMenu(update, bot)
	}
}

func punishMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err := services.GetModelDataWhere(utils.GroupInfo.GroupId, &punishment)

	var btns [][]model.ButtonInfo
	utils.Json2Button2("punish.json", &btns)

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
	case "ban":
		punishment.PunishType = model.PunishTypeBan
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
	punishment.WarningAfterPunish = text

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
	}
	//todo 根据class类型分别处理
	actionMsg := "警告 "
	if punishment.PunishType == model.PunishTypeBan {
		actionMsg = "禁言"
	} else if punishment.PunishType == model.PunishTypeKick {
		actionMsg = "踢出"
	} else if punishment.PunishType == model.PunishTypeBanAndKick {
		actionMsg = "踢出+禁言"
	} else if punishment.PunishType == model.PunishTypeRevoke {
		actionMsg = "仅撤回消息+不惩罚"
	} else if punishment.PunishType == model.PunishTypeWarning {
		afterMsg := "禁言"
		if punishment.WarningAfterPunish == "kick" {
			afterMsg = "踢出"
		} else if punishment.WarningAfterPunish == "banAndKick" {
			afterMsg = "踢出+禁言"
		}
		actionMsg = fmt.Sprintf("警告%d次后 %s", punishment.WarningCount, afterMsg)
	}

	content = content + actionMsg
	switch class {
	case "spam":
		spamsSetting.Punishment = punishment
		updateSpamMsg()
	case "prohibited":
		prohibitedSetting.Punishment = punishment
		updateProhibitedSettingMsg()
	}
	return content
}

func banTimeMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("🔇 违禁词\n\n当前设置：%d分钟 \n👉 输入处罚禁言的时长（分钟，例如：60）：", punishment.BanTime)
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
	punishment.BanTime = time
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

	//updateMsg()
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
		backAction = "flood_setting"
	} else if class == "spam" {
		backAction = "spam_setting"
	} else if class == "prohibited" {
		backAction = "prohibited_setting"
	}
	return backAction
}
