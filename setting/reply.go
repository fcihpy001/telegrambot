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

var replySetting model.ReplySetting
var reply model.Reply
var replySelect model.SelectInfo

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

	} else if cmd == "auto_reply_keyword_add" {
		addKeywordMenu(update, bot)
	} else if cmd == "auto_reply_keyword_delete" {
		deleteKeywordMenu(update, bot)
	} else if cmd == "auto_reply_delete_time" {
		deleteReplyTimeHandler(update, bot, params)
	} else if cmd == "auto_reply_keyword_trigger_type" {
		triggerMatchType(update, bot, params)
	}
}

func replyMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	replySetting, _ = services.GetReplySetting(utils.GroupInfo.GroupId)
	replySetting.ChatId = utils.GroupInfo.GroupId
	updateSelectInfo()

	var buttons [][]model.ButtonInfo
	utils.Json2Button2("./config/reply.json", &buttons)
	fmt.Println(&buttons)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btnArr := buttons[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArr); j++ {
			btn := btnArr[j]
			if btn.Text == "启用" && replySetting.Enable {
				btn.Text = "✅启用"
			} else if btn.Text == "关闭" && !replySetting.Enable {
				btn.Text = "✅关闭"
			}
			if strings.HasPrefix(btn.Data, "auto_reply_delete_time") && btn.Data == fmt.Sprintf("auto_reply_delete_time:%d", replySetting.DeleteReplyTime) {
				btn.Text = fmt.Sprintf("✅%d", replySetting.DeleteReplyTime)

			} else if btn.Text == "否" && replySetting.DeleteReplyTime == 0 {
				btn.Text = "✅否"
			}
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard_enable := utils.MakeKeyboard(rows)

	rows_disable := append(rows[:1], rows[4:]...)
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

// 状态处理
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

// 自动删除时间
func deleteReplyTimeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	index := strings.Split(params, "&")
	text := index[0]
	col := index[1]
	colInt, _ := strconv.Atoi(col)

	time, _ := strconv.Atoi(params)
	replySetting.DeleteReplyTime = time

	//	取消原来的选中
	utils.ReplEnableyMenuMarkup.InlineKeyboard[replySelect.Row][replySelect.Column].Text = replySelect.Text
	//	选中新的
	utils.ReplEnableyMenuMarkup.InlineKeyboard[2][colInt].Text = "✅" + text
	//	更新选中
	replySelect.Text = text
	replySelect.Row = 2
	replySelect.Column = colInt

	content := updateReplySettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.ReplEnableyMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

}

// 添加关键词: step1-给出提示
func addKeywordMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//遍历着急词列表,列出每个关键词，按行拼接
	text := ""
	for _, reply := range replySetting.ReplyList {
		text = text + reply.KeyWorld + "\n"
	}
	content := fmt.Sprintf("💬 关键词回复\n\n已添加的关键词：\n%s\n\n👉第一步 请输入关键词：", text)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "请输入关键词",
		Selective:             false,
	}
	bot.Send(msg)
}

// 添加关键词: step2-收到关键词输入回复，提示用户输入回复内容
func AddKeywordResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//删除前一个消息
	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	bot.Send(deleteMsg)

	//组装model数据
	reply.KeyWorld = update.Message.Text
	reply.ChatId = utils.GroupInfo.GroupId
	reply.ReplySettingID = replySetting.ID

	//发送提示消息
	content := fmt.Sprintf("💬 关键词回复\n\n👉 第二步 请输入关键词%s的回复内容（支持图片，表情，视频，文件等）：", update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "请输入回复内容",
		Selective:             false,
	}
	bot.Send(msg)
}

// 添加关键词: step3-收到回复内容，提示用户选择回复触发方式
func AddKeywordReplyResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//删除前一个消息
	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	bot.Send(deleteMsg)

	reply.ReplyWorld = update.Message.Text
	btn1 := model.ButtonInfo{
		Text:    "精准匹配",
		Data:    "auto_reply_keyword_trigger_type:精准匹配",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "模糊匹配",
		Data:    "auto_reply_keyword_trigger_type:模糊匹配",
		BtnType: model.BtnTypeData,
	}
	row := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row}
	keyboard := utils.MakeKeyboard(rows)

	content := "💬 关键词回复\n\n👉最后，请选择回复触发方式："
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// 添加关键词: step4-收到回复触发方式，提示用户添加成功，更新model数据
func triggerMatchType(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	if params == "精准匹配" {
		reply.MatchAll = true
	} else {
		reply.MatchAll = false
	}
	replySetting.ReplyList = append(replySetting.ReplyList, reply)
	updateReplySettingMsg()
	content := "✅操作成功"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "auto_reply_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 删除关键词
func deleteKeywordMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//遍历着急词列表,列出每个关键词，按行拼接
	text := ""
	for _, reply := range replySetting.ReplyList {
		text = text + reply.KeyWorld + "\n"
	}
	content := fmt.Sprintf("💬请输入要删除的关键词，一次只能删除一个，回复关键词名：\n                \n已添加的关键词：\n%s", text)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "tetetet",
		Selective:             true,
	}
	bot.Send(msg)
}

// 删除关键词回应处理
func DeleteKeywordResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	text := update.Message.Text
	//从列表中找到这个关键词，然后删除
	for i, v := range replySetting.ReplyList {
		if v.KeyWorld == text {
			replySetting.ReplyList = append(replySetting.ReplyList[:i], replySetting.ReplyList[i+1:]...)
			//删除数据库中的数据
			_ = services.DeleteReply(text, utils.GroupInfo.GroupId)
			break
		}
	}
	content := "✅操作成功"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "auto_reply_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateReplySettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 更新model数据，并将信息入库
func updateReplySettingMsg() string {
	content := "💬 关键词回复\n\n在群组中使用命令：\n/filter 添加自动回复规则\n/stop 删除自动回复规则\n/filters 所有自动回复规则列表\n查看命令帮助\n\n已添加的关键词：\n"
	if replySetting.Enable == false {
		content = "💬 关键词回复\n\n当前状态：关闭❌"
		return content
	}
	for _, v := range replySetting.ReplyList {
		if v.MatchAll {
			content = content + "\n- " + v.KeyWorld
		} else {
			content = content + "\n* " + v.KeyWorld
		}
	}
	content = content + "\n" + "\n- 表示精准触发\n * 表示包含触发"

	services.SaveModel(&replySetting, replySetting.ChatId)
	reply = model.Reply{}
	return content
}

// 回复逻辑处理
func HandlerAutoReply(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	//读取配置状态
	replySetting, _ = services.GetReplySetting(utils.GroupInfo.GroupId)
	if replySetting.Enable == false {
		return false
	}

	// 获取用户发送的消息文本
	messageText := update.Message.Text

	//从数据库中取取出所有的自动回复词库
	relyList, err := services.GetAllReply(utils.GroupInfo.GroupId)
	if err != nil {
		log.Println(err)
	}
	//根据收到消息，与ll中每个Model的keyworkd比较，如果matchAll为true，那么就是完全匹配，否则就是包含匹配
	for _, v := range relyList {
		if v.MatchAll {
			if messageText == v.KeyWorld {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, v.ReplyWorld)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				return true
			}
		} else {
			if strings.Contains(messageText, v.KeyWorld) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, v.ReplyWorld)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				return true
			}
		}
	}
	return false
}

func updateSelectInfo() {
	str := "否"
	col := 0
	if replySetting.DeleteReplyTime == 1 {
		str = "1"
		col = 1
	} else if replySetting.DeleteReplyTime == 5 {
		str = "5"
		col = 2
	} else if replySetting.DeleteReplyTime == 10 {
		str = "10"
		col = 3
	} else if replySetting.DeleteReplyTime == 30 {
		str = "30"
		col = 4
	}
	replySelect = model.SelectInfo{
		Text:   str,
		Row:    2,
		Column: col,
	}
}
