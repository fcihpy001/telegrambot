package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
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

	} else if cmd == "dark_model_status" {
		statusHandler(update, bot, params)

	}
}

func inviteSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		inviteSetting = services.GetInviteSettings(update.CallbackQuery.Message.Chat.ID)
	}
	_ = services.GetModelData(utils.GroupInfo.GroupId, &inviteSetting)
	inviteSetting.ChatId = utils.GroupInfo.GroupId

	//btn12txt := "启用"
	//btn13txt := "✅关闭"
	//if inviteSetting.Enable {
	//	btn12txt = "✅启用"
	//	btn13txt = "关闭"
	//}
	//
	//btn22txt := "启用"
	//btn23txt := "✅关闭"
	//if inviteSetting.Enable {
	//	btn12txt = "✅启用"
	//	btn13txt = "关闭"
	//}
	//
	//btn32txt := "启用"
	//btn33txt := "✅关闭"
	//if inviteSetting.Enable {
	//	btn12txt = "✅启用"
	//	btn13txt = "关闭"
	//}
	//
	//btn11 := model.ButtonInfo{
	//	Text:    "是否启用",
	//	Data:    "toast",
	//	BtnType: model.BtnTypeData,
	//}
	//btn12 := model.ButtonInfo{
	//	Text:    btn12txt,
	//	Data:    "group_invite_status_enable",
	//	BtnType: model.BtnTypeData,
	//}
	//btn13 := model.ButtonInfo{
	//	Text:    btn13txt,
	//	Data:    "group_invite_status_disable",
	//	BtnType: model.BtnTypeData,
	//}
	//btn21 := model.ButtonInfo{
	//	Text:    "入群自动生成",
	//	Data:    "toast",
	//	BtnType: model.BtnTypeData,
	//}
	//btn22 := model.ButtonInfo{
	//	Text:    btn22txt,
	//	Data:    "group_invite_autogenerate_enable",
	//	BtnType: model.BtnTypeData,
	//}
	//btn23 := model.ButtonInfo{
	//	Text:    btn23txt,
	//	Data:    "group_invite_autogenerate_disable",
	//	BtnType: model.BtnTypeData,
	//}
	//btn31 := model.ButtonInfo{
	//	Text:    "邀请成功提醒",
	//	Data:    "toast",
	//	BtnType: model.BtnTypeData,
	//}
	//btn32 := model.ButtonInfo{
	//	Text:    btn32txt,
	//	Data:    "group_invite_notify_enable",
	//	BtnType: model.BtnTypeData,
	//}
	//btn33 := model.ButtonInfo{
	//	Text:    btn33txt,
	//	Data:    "group_invite_notify_disable",
	//	BtnType: model.BtnTypeData,
	//}
	//
	//btn4 := model.ButtonInfo{
	//	Text:    "🦚设置链接过期时间",
	//	Data:    "group_inivte_expire_time",
	//	BtnType: model.BtnTypeData,
	//}
	//btn5 := model.ButtonInfo{
	//	Text:    "🍇设置单链接取大邀请数",
	//	Data:    "group_inivte_people_limit",
	//	BtnType: model.BtnTypeData,
	//}
	//btn6 := model.ButtonInfo{
	//	Text:    "🍵设置链接生成数量上限",
	//	Data:    "group_inivte_link_limit",
	//	BtnType: model.BtnTypeData,
	//}
	//btn71 := model.ButtonInfo{
	//	Text:    "导出邀请统计",
	//	Data:    "group_inivte_export",
	//	BtnType: model.BtnTypeData,
	//}
	//btn72 := model.ButtonInfo{
	//	Text:    "清空所有邀请统计",
	//	Data:    "group_inivte_clear",
	//	BtnType: model.BtnTypeData,
	//}
	//btn81 := model.ButtonInfo{
	//	Text:    "🏠返回",
	//	Data:    "go_setting",
	//	BtnType: model.BtnTypeData,
	//}
	//row1 := []model.ButtonInfo{btn11, btn12, btn13}
	//row2 := []model.ButtonInfo{btn21, btn22, btn23}
	//row3 := []model.ButtonInfo{btn31, btn32, btn33}
	//row4 := []model.ButtonInfo{btn4}
	//row5 := []model.ButtonInfo{btn5}
	//row6 := []model.ButtonInfo{btn6}
	//row7 := []model.ButtonInfo{btn71, btn72}
	//row8 := []model.ButtonInfo{btn81}
	//rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5, row6, row7, row8}

	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/invite.json", &btns)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			btn := btnArray[j]
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

//func group_invite_status(update *tgbotapi.Update, enable bool) {
//	if enable {
//		utils.InviteMenuMarkup.InlineKeyboard[0][1].Text = "✅启用"
//		utils.InviteMenuMarkup.InlineKeyboard[0][2].Text = "关闭"
//		inviteSetting.Enable = true
//	} else {
//		utils.InviteMenuMarkup.InlineKeyboard[0][1].Text = "启用"
//		utils.InviteMenuMarkup.InlineKeyboard[0][2].Text = "✅关闭"
//		inviteSetting.Enable = false
//	}
//
//	content := updateInviteSettingMsg()
//	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.InviteMenuMarkup)
//	_, err := mgr.bot.Send(msg)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func (mgr *group.GroupManager) group_invite_autogenerate(update *tgbotapi.Update, enable bool) {
//	if enable {
//		utils.InviteMenuMarkup.InlineKeyboard[1][1].Text = "✅启用"
//		utils.InviteMenuMarkup.InlineKeyboard[1][2].Text = "关闭"
//		inviteSetting.Enable = true
//	} else {
//		utils.InviteMenuMarkup.InlineKeyboard[1][1].Text = "启用"
//		utils.InviteMenuMarkup.InlineKeyboard[1][2].Text = "✅关闭"
//		inviteSetting.Enable = false
//	}
//
//	content := updateInviteSettingMsg()
//	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.InviteMenuMarkup)
//	_, err := mgr.bot.Send(msg)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func (mgr *group.GroupManager) group_invite_notify(update *tgbotapi.Update, enable bool) {
//	if enable {
//		utils.InviteMenuMarkup.InlineKeyboard[2][1].Text = "✅启用"
//		utils.InviteMenuMarkup.InlineKeyboard[2][2].Text = "关闭"
//		inviteSetting.Enable = true
//	} else {
//		utils.InviteMenuMarkup.InlineKeyboard[2][1].Text = "启用"
//		utils.InviteMenuMarkup.InlineKeyboard[2][2].Text = "✅关闭"
//		inviteSetting.Enable = false
//	}
//
//	content := updateInviteSettingMsg()
//	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, utils.InviteMenuMarkup)
//	_, err := mgr.bot.Send(msg)
//	if err != nil {
//		log.Println(err)
//	}
//}

func updateInviteSettingMsg() string {
	content := "🔗 【toplink官方群】邀请链接生成\n\n开启后群组中成员使用 /link 指令自动生成链接/查询邀请统计\n\n防作弊：\n└ 只有第一次进群视为有效邀请数，退群再用其他人的链接加群不计算邀请数\n"
	enableMsg := "┌当前状态：关闭 ❌"
	if inviteSetting.Enable {
		enableMsg = "当前状态：开启 ✅"
	}
	InviteCount := "├总邀请人数：" + string(rune(inviteSetting.InviteCount)) + "\n"
	linkeExpireTime := "├邀请链接有效期：不限制 \n"
	if inviteSetting.LinkExpireTime > 0 {
		linkeExpireTime = "├邀请链接有效期：" + string(rune(inviteSetting.LinkExpireTime)) + "天\n"
	}

	InvitePeopleLimit := "├最大邀请人数：无限制\n"
	if inviteSetting.LinkExpireTime > 0 {
		InvitePeopleLimit = "├最大邀请人数：" + string(rune(inviteSetting.InvitePeopleLimit)) + "\n"
	}

	InviteLinkLimit := "└生成数量上限： 无限制     已生成数量：0\n"
	if inviteSetting.LinkExpireTime > 0 {
		InviteLinkLimit = "└生成数量上限： " + string(rune(inviteSetting.InviteLinkLimit)) + "已生成数量：0\n"
	}

	content = content + enableMsg + "\n" + InviteCount + linkeExpireTime + InvitePeopleLimit + InviteLinkLimit
	services.SaveInviteSettings(&inviteSetting)
	return content
}
