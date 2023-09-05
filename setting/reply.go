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

var replySetting model.ReplySetting

func ReplyHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "auto_reply_menu" {
		replyMenu(update, bot)

	} else if cmd == "auto_reply_status" {
		replyStatusHandler(update, bot, params == "enable")

	}
}

func replyMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	kkk := model.Reply{
		ChatId:     999,
		KeyWorld:   "hello",
		ReplyWorld: "How are you",
	}
	replySetting := services.GetReplySettings(update.CallbackQuery.Message.Chat.ID)
	replySetting.ChatId = update.CallbackQuery.Message.Chat.ID
	replySetting.KeywordReply = append(replySetting.KeywordReply, kkk)
	services.SaveReplySettings(&replySetting)
	fmt.Println("reply_data:", replySetting)
	btn12txt := "启用"
	btn13txt := "✅关闭"
	if replySetting.Enable {
		btn12txt = "✅启用"
		btn13txt = "关闭"
	}

	btn31txt := "✅否"
	btn32txt := "1"
	btn33txt := "5"
	btn34txt := "10"
	btn35txt := "30"

	if replySetting.DeleteReplyTime == 1 {
		btn31txt = "否"
		btn32txt = "✅1"
		btn33txt = "5"
		btn34txt = "10"
		btn35txt = "30"
	} else if replySetting.DeleteReplyTime == 5 {
		btn31txt = "否"
		btn32txt = "1"
		btn33txt = "✅5"
		btn34txt = "10"
		btn35txt = "30"
	} else if replySetting.DeleteReplyTime == 10 {
		btn31txt = "否"
		btn32txt = "1"
		btn33txt = "5"
		btn34txt = "✅10"
		btn35txt = "30"
	}
	btn11 := model.ButtonInfo{
		Text:    "是否启用",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    btn12txt,
		Data:    "reply_status_enable",
		BtnType: model.BtnTypeData,
	}
	btn13 := model.ButtonInfo{
		Text:    btn13txt,
		Data:    "reply_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "自动删除回复消息(分钟)↘️",
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}

	btn31 := model.ButtonInfo{
		Text:    btn31txt,
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn32 := model.ButtonInfo{
		Text:    btn32txt,
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn33 := model.ButtonInfo{
		Text:    btn33txt,
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn34 := model.ButtonInfo{
		Text:    btn34txt,
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn35 := model.ButtonInfo{
		Text:    btn35txt,
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "➕添加关键词",
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "🚽删除关键词",
		Data:    "group_invite_status_disable",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12, btn13}

	row2 := []model.ButtonInfo{btn21}
	row3 := []model.ButtonInfo{btn31, btn32, btn33, btn34, btn35}
	row4 := []model.ButtonInfo{btn41, btn42}
	row5 := []model.ButtonInfo{btn51}
	rows_enable := [][]model.ButtonInfo{row1, row2, row3, row4, row5}
	rows_disable := [][]model.ButtonInfo{row1, row5}

	keyboard_enable := utils.MakeKeyboard(rows_enable)
	keyboard_disable := utils.MakeKeyboard(rows_disable)

	utils.ReplEnableyMenuMarkup = keyboard_enable
	utils.ReplDisableMenuMarkup = keyboard_disable

	var keyboard tgbotapi.InlineKeyboardMarkup
	if replySetting.Enable {
		keyboard = keyboard_enable
	} else {
		keyboard = keyboard_disable
	}

	//要读取用户设置的数据
	content := updateReplySettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func replyStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {

	replySetting.Enable = enable
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	if enable {
		utils.ReplEnableyMenuMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.ReplEnableyMenuMarkup.InlineKeyboard[0][2].Text = "关闭"
		keyboard = utils.ReplEnableyMenuMarkup
	} else {
		utils.ReplDisableMenuMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.ReplDisableMenuMarkup.InlineKeyboard[0][2].Text = "✅关闭"
		keyboard = utils.ReplDisableMenuMarkup
	}

	content := updateReplySettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateReplySettingMsg() string {
	content := "💬 关键词回复\n\n在群组中使用命令：\n/filter 添加自动回复规则\n/stop 删除自动回复规则\n/filters 所有自动回复规则列表\n查看命令帮助\n\n已添加的关键词：\n"
	if replySetting.Enable == false {
		content = "💬 关键词回复\n\n当前状态：关闭❌"
		return content
	}
	fmt.Println("reply_keyworld", replySetting.KeywordReply)

	enableMsg := "* match world"

	content = content + enableMsg + "\n" + "\n- 表示精准触发\n * 表示包含触发"

	services.SaveReplySettings(&replySetting)
	return content
}
