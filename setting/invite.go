package setting

import (
	"encoding/json"
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

var inviteSetting model.InviteSetting

func InviteHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "invite_setting_menu" {
		inviteSettingMenu(update, bot)

	} else if cmd == "invite_setting_status" {
		inviteStatusHandler(update, bot, params == "enable")

	} else if cmd == "invite_setting_autogenerate" {
		linkGenerateHandler(update, bot, params == "enable")

	} else if cmd == "invite_setting_notify" {
		inviteNotifyHandler(update, bot, params == "enable")

	} else if cmd == "invite_setting_expire_time" {
		inviteExpireTimeMenu(update, bot)

	} else if cmd == "invite_setting_limit_people" {
		invitePeopleMenu(update, bot)

	} else if cmd == "invite_setting_limit_link" {
		inviteLinkLimitMenu(update, bot)

	} else if cmd == "invite_setting_export" {
		exportInviteData(update, bot)

	} else if cmd == "invite_setting_clear" {
		clearInviteMenu(update, bot)

	} else if cmd == "invite_setting_clear_confirm" {
		clearInviteData(update, bot)

	}
}

func inviteSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_ = services.GetModelData(utils.GroupInfo.GroupId, &inviteSetting)
	inviteSetting.ChatId = utils.GroupInfo.GroupId

	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/invite.json", &btns)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			btn := btnArray[j]
			updateInviteButtonStatus(&btn)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.InviteMenuMarkup = keyboard

	//要读取用户设置的数据
	content := updateInviteSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 状态处理
func inviteStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
	inviteSetting.Enable = enable
	if enable {
		utils.InviteMenuMarkup.InlineKeyboard[0][1].Text = "✅启用"
		utils.InviteMenuMarkup.InlineKeyboard[0][2].Text = "关闭"
	} else {
		utils.InviteMenuMarkup.InlineKeyboard[0][1].Text = "启用"
		utils.InviteMenuMarkup.InlineKeyboard[0][2].Text = "✅关闭"
	}

	content := updateInviteSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.InviteMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 自动生成链接
func linkGenerateHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
	inviteSetting.AutoGenerate = enable
	if enable {
		utils.InviteMenuMarkup.InlineKeyboard[1][1].Text = "✅启用"
		utils.InviteMenuMarkup.InlineKeyboard[1][2].Text = "关闭"
	} else {
		utils.InviteMenuMarkup.InlineKeyboard[1][1].Text = "启用"
		utils.InviteMenuMarkup.InlineKeyboard[1][2].Text = "✅关闭"
	}

	content := updateInviteSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.InviteMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 邀请成功是否通知
func inviteNotifyHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, enable bool) {
	inviteSetting.Notify = enable
	if enable {
		utils.InviteMenuMarkup.InlineKeyboard[2][1].Text = "✅启用"
		utils.InviteMenuMarkup.InlineKeyboard[2][2].Text = "关闭"
	} else {
		utils.InviteMenuMarkup.InlineKeyboard[2][1].Text = "启用"
		utils.InviteMenuMarkup.InlineKeyboard[2][2].Text = "✅关闭"
	}

	content := updateInviteSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.InviteMenuMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 邀请链接有效期
func inviteExpireTimeMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("请回复链接过期时间(不限制请输入：0)\n\n注意：此设置仅应用在新生成的链接中，不会修改已生成的链接\n\n格式：年-月-日 时:分\n例如：%s (点击复制)", utils.CurrentTime())
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

func InviteExpireTimeResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	inviteSetting.ExpireDate = update.Message.Text
	content := "✅设置成功，点击按钮返回"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "invite_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateInviteSettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 邀请人数上限
func invitePeopleMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("邀请达到设定人数后链接失效\n\n注意：此设置仅应用在新生成的链接中，不会修改已生成的链接\n\n请回复单个链接最大邀请人数(不限制请输入：0)")
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

func InvitePeopleLimitResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	count, err := strconv.Atoi(update.Message.Text)

	inviteSetting.MemberLimit = count
	content := "✅设置成功，点击按钮返回"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "invite_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateInviteSettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 生成链接数量上限
func inviteLinkLimitMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("生成链接数量达到设定数量后，不再生成新的链接\n\n注意：此设置仅应用在新生成的链接中，不会修改已生成的链接\n\n请回复生成链接数量上限(不限制请输入：0)")
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	keybord := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("返回"),
		))

	msg.ReplyMarkup = keybord
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "请输入生成链接数量上限",
	}
	bot.Send(msg)
}

func InviteLinkLimitResult(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	count, err := strconv.Atoi(update.Message.Text)

	inviteSetting.InviteLinkLimit = count
	content := "✅设置成功，点击按钮返回"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "invite_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	updateInviteSettingMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 导出邀请数据
func exportInviteData(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := "正在导出邀请数据，请稍后..."
	if inviteSetting.InviteCount == 0 {
		content = "没有邀请数据，无需导出"
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, content)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 清除邀请数据
func clearInviteMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := "🚨🚨 请注意，即将清空所有邀请链接和邀请数据，操作不可恢复，是否继续："
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "flood_setting_menu",
		BtnType: model.BtnTypeData,
	}
	btn2 := model.ButtonInfo{
		Text:    "❗️确认删除",
		Data:    "invite_setting_clear_confirm",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1, btn2}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	bot.Send(msg)
}

func clearInviteData(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	inviteSetting = model.InviteSetting{}
	err = services.DeleteInviteData()
	if err != nil {
		log.Println(err)
	}
	content := "邀请链接和邀请数据已清空，点击按钮返回"
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "flood_setting_menu",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	bot.Send(msg)
}

func updateInviteSettingMsg() string {
	content := fmt.Sprintf("🔗 【%s】邀请链接生成\n\n开启后群组中成员使用"+
		" /link 指令自动生成链接/查询邀请统计\n\n防作弊：\n"+
		"└ 只有第一次进群视为有效邀请数，退群再用其他人的链接加群不计算邀请数\n", utils.GroupInfo.GroupName)
	enableMsg := "┌当前状态：关闭 ❌"
	if inviteSetting.Enable {
		enableMsg = "当前状态：开启 ✅"
	}

	InviteCount := "├总邀请人数：" + strconv.Itoa(inviteSetting.InviteCount) + "\n"
	linkExpireTime := "├邀请链接有效期：不限制 \n"
	if inviteSetting.ExpireDate != "0" {
		linkExpireTime = "├邀请链接有效期：" + inviteSetting.ExpireDate + "\n"
	}

	InvitePeopleLimit := "├最大邀请人数：无限制\n"
	if inviteSetting.MemberLimit > 0 {
		InvitePeopleLimit = "├最大邀请人数：" + strconv.Itoa(inviteSetting.MemberLimit) + "\n"
	}

	InviteLinkLimit := "└生成数量上限： 无限制     已生成数量：0\n"
	if inviteSetting.InviteLinkLimit > 0 {
		InviteLinkLimit = "└生成数量上限： " + string(rune(inviteSetting.InviteLinkLimit)) + "已生成数量：0\n"
	}

	content = content + enableMsg + "\n" + InviteCount + linkExpireTime + InvitePeopleLimit + InviteLinkLimit
	services.SaveModel(&inviteSetting, inviteSetting.ChatId)
	return content
}

func updateInviteButtonStatus(btn *model.ButtonInfo) {
	if btn.Data == "invite_setting_status:enable" && inviteSetting.Enable {
		btn.Text = "✅启用"
	} else if btn.Data == "invite_setting_status:disable" && !inviteSetting.Enable {
		btn.Text = "✅关闭"
	} else if btn.Data == "invite_setting_autogenerate:enable" && inviteSetting.AutoGenerate {
		btn.Text = "✅启用"
	} else if btn.Data == "invite_setting_autogenerate:disable" && !inviteSetting.AutoGenerate {
		btn.Text = "✅关闭"
	} else if btn.Data == "invite_setting_notify:enable" && inviteSetting.Notify {
		btn.Text = "✅通知"
	} else if btn.Data == "invite_setting_notify:disable" && !inviteSetting.Notify {
		btn.Text = "✅不通知"
	}
}

func GetInviteLink(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//获取群的邀请配置信息
	inviteSetting := model.InviteSetting{}
	_ = services.GetModelData(update.Message.Chat.ID, &inviteSetting)
	if inviteSetting.Enable == false {
		utils.SendText(update.Message.Chat.ID, "管理员关闭了生成邀请链接的开关", bot)
		return
	}
	invite := model.InviteRecord{}
	where := fmt.Sprintf("chat_id = %d and uid = %d", update.Message.Chat.ID, update.Message.From.ID)
	_ = services.GetModelWhere(where, &invite)
	link := invite.InviteLink
	if invite.ID > 0 {
		msg := fmt.Sprintf("🔗 @%s 您的专属链接:\n %s (点击复制)\n\n👉 👉 当前总共邀请%d人\n\n", update.Message.From.FirstName, link, invite.InviteCount)
		utils.SendText(update.Message.Chat.ID, msg, bot)
		return
	}

	expireDate := int(time.Now().Unix() + 86400*365)
	if len(inviteSetting.ExpireDate) > 0 {
		expireDateStr, _ := time.Parse("2006-01-02 15:04", inviteSetting.ExpireDate)
		expireDate = int(expireDateStr.Unix())
	}

	memberLimit := 9999
	if inviteSetting.MemberLimit > 0 {
		memberLimit = inviteSetting.MemberLimit
	}

	linkName := update.Message.From.FirstName + "的专属链接"
	config := tgbotapi.CreateChatInviteLinkConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: update.Message.Chat.ID,
		},
		Name:               linkName,
		ExpireDate:         expireDate,
		MemberLimit:        memberLimit,
		CreatesJoinRequest: false,
	}
	//tip CreatesJoinRequest与MemberLimit不能同时设置

	resp, err := bot.Request(config)
	if err != nil {
		fmt.Println("invite get failed:", err)
		utils.SendText(update.Message.Chat.ID, "机器人的管理权限不足", bot)
		return
	}
	m := map[string]interface{}{}

	json.Unmarshal(resp.Result, &m)
	link = m["invite_link"].(string)

	msg := fmt.Sprintf("🔗 @%s 您的专属链接:\n %s (点击复制)\n\n👉 👉 当前总共邀请0人\n\n（本消息5分钟自毁）", update.Message.From.FirstName, link)
	utils.SendText(update.Message.Chat.ID, msg, bot)

	//	将邀请链接入库
	invite.InviteLink = link
	invite.ChatId = update.Message.Chat.ID
	invite.Uid = update.Message.From.ID
	invite.LinkName = linkName
	invite.ChatName = update.Message.Chat.Title
	invite.ChatType = update.Message.Chat.Type
	invite.ExpireDate = expireDate
	invite.MemberLimit = inviteSetting.MemberLimit
	invite.CreatesJoinRequest = false
	services.SaveModel(&invite, invite.ChatId)
}

func UpdateInviteRecord(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	link := update.ChatJoinRequest.InviteLink.InviteLink
	count := update.ChatJoinRequest.InviteLink.PendingJoinRequestCount
	userid := update.ChatJoinRequest.From.ID
	record := model.InviteRecord{}
	record.InviteLink = link
	record.LinkName = update.ChatJoinRequest.InviteLink.Name
	record.ChatId = update.ChatJoinRequest.Chat.ID
	record.Uid = userid
	record.InviteCount = count
	services.SaveModel(&record, record.ChatId)

}
