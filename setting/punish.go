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

func punishMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err := services.GetModelDataWhere(update.CallbackQuery.Message.Chat.ID, &punishment)

	punishType := []string{"警告", "禁言", "踢出", "踢出+封禁", "仅撤回消息+不惩罚"}
	punishTypeEn := []string{"warning", "ban", "kick", "banAndKick", "revoke"}
	warningCounts := []string{"1", "2", "3", "4", "5"}
	warningAfterActions := []string{"禁言", "踢出", "踢出+封禁"}
	warningAfterActionsEn := []string{"ban", "kick", "banAndKick"}

	//惩罚方法
	typeRow := []model.ButtonInfo{}
	rows := [][]model.ButtonInfo{}
	rows2 := [][]model.ButtonInfo{}
	rows3 := [][]model.ButtonInfo{}
	for i := 0; i < 5; i++ {
		btn := model.ButtonInfo{
			Text:    punishType[i],
			Data:    "punish_setting_type:" + punishTypeEn[i],
			BtnType: model.BtnTypeData,
		}
		typeRow = append(typeRow, btn)
		if i == 2 {
			rows = append(rows, typeRow)
			typeRow = []model.ButtonInfo{}
		}
	}
	rows = append(rows, typeRow)
	rows2 = rows
	rows3 = rows

	//警告次数
	tip1 := model.ButtonInfo{
		Text:    "警告次数",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	tip1Row := []model.ButtonInfo{tip1}
	warningRow := []model.ButtonInfo{}
	for i := 0; i < len(warningCounts); i++ {
		btn := model.ButtonInfo{
			Text:    warningCounts[i],
			Data:    "punish_setting_count:" + warningCounts[i],
			BtnType: model.BtnTypeData,
		}
		warningRow = append(warningRow, btn)
	}
	rows = append(rows, tip1Row)
	rows = append(rows, warningRow)

	//达到警告次数后动作
	tip2 := model.ButtonInfo{
		Text:    "达到警告次数后动作",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	tip2Row := []model.ButtonInfo{tip2}
	actionRow := []model.ButtonInfo{}
	for i := 0; i < len(warningAfterActions); i++ {
		btn := model.ButtonInfo{
			Text:    warningAfterActions[i],
			Data:    "punish_setting_action:" + warningAfterActionsEn[i],
			BtnType: model.BtnTypeData,
		}
		actionRow = append(actionRow, btn)
	}
	rows = append(rows, tip2Row)
	rows = append(rows, actionRow)

	backMenu := ""
	if class == "flood" {
		backMenu = "flood_setting"
	} else if class == "spam" {
		backMenu = "spam_setting"
	}

	timeRow := model.ButtonInfo{
		Text:    "设置禁言时长",
		Data:    "punish_setting_time",
		BtnType: model.BtnTypeData,
	}
	timeRows := []model.ButtonInfo{timeRow}
	rows2 = append(rows2, timeRows)

	rows = append(rows, timeRows)
	btn71 := model.ButtonInfo{
		Text:    "返回",
		Data:    backMenu,
		BtnType: model.BtnTypeData,
	}
	backRows := []model.ButtonInfo{btn71}
	rows = append(rows, backRows)
	keyboard := utils.MakeKeyboard(rows)
	utils.PunishMenuMarkup = keyboard

	//禁言键盘  类型+时长
	rows2 = append(rows2, backRows)
	keyboard2 := utils.MakeKeyboard(rows2)
	utils.PunishMenuMarkup2 = keyboard2

	//仅动作键盘
	//type3_rows := typeRows
	rows3 = append(rows3, backRows)
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

func PunishSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "punish_setting_class" {
		class = params
		punishMenu(update, bot)

	} else if cmd == "punish_setting_type" {
		punishTypeHandler(update, bot, params)

	} else if cmd == "punish_setting_count" {
		count, _ := strconv.Atoi(params)
		warningCountHandler(update, bot, count)

	} else if cmd == "punish_setting_action" {
		warningActionHandler(update, bot, params)
	} else if cmd == "punish_setting_time" {

	}

}

func punishTypeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	switch params {
	case "warning":
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
func warningActionHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if params == "ban" {
		utils.PunishMenuMarkup.InlineKeyboard[5][0].Text = "✅禁言"
		utils.PunishMenuMarkup.InlineKeyboard[5][1].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[5][2].Text = "踢出+封禁"
		punishment.WarningAfterPunish = model.PunishTypeBan

	} else if params == "kick" {
		utils.PunishMenuMarkup.InlineKeyboard[5][0].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[5][1].Text = "✅踢出"
		utils.PunishMenuMarkup.InlineKeyboard[5][2].Text = "踢出+封禁"
		punishment.WarningAfterPunish = model.PunishTypeKick

	} else if params == "banAndKick" {
		utils.PunishMenuMarkup.InlineKeyboard[5][0].Text = "禁言"
		utils.PunishMenuMarkup.InlineKeyboard[5][1].Text = "踢出"
		utils.PunishMenuMarkup.InlineKeyboard[5][2].Text = "✅踢出+封禁"
		punishment.WarningAfterPunish = model.PunishTypeBanAndKick
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
	actionMsg := "警告"

	if punishment.PunishType == model.PunishTypeBan {
		actionMsg = "禁言"
	} else if punishment.PunishType == model.PunishTypeKick {
		actionMsg = "踢出"
	} else if punishment.PunishType == model.PunishTypeBanAndKick {
		actionMsg = "踢出+禁言"
	} else if punishment.PunishType == model.PunishTypeRevoke {
		actionMsg = "仅撤回消息+不惩罚"
	} else if punishment.PunishType == model.PunishTypeWarning {
		actionMsg = fmt.Sprintf("警告%d次后 %s", prohibitedSetting.WarningCount, actionMap[prohibitedSetting.WarningAfterPunish])
	}

	content = content + actionMsg
	switch class {
	case "spam":
		spamsSetting.Punishment = punishment
		updateSpamMsg()
	case "prohibit":
		prohibitedSetting.Punishment = punishment
		updateProhibitedSettingMsg()

	}
	return content
}
