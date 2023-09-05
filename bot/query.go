package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/group"
	"telegramBot/lucky"
	"telegramBot/setting"
	"telegramBot/utils"
)

// 处理行内按钮事件
func (bot *SmartBot) handleQuery(update *tgbotapi.Update) {
	query := update.CallbackQuery.Data
	fmt.Println("query command--", query)

	if strings.HasPrefix(query, "lucky") {
		lucky.LuckyHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "group") { //群组设置
		group.GroupHandlerQuery(update, bot.bot)

	} else if strings.HasPrefix(query, "settings") { //设置
		setting.Settings(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Chat.Type, "", bot.bot)

	} else if query == "go_setting" { //设置主菜单
		bot.go_setting(update)

	} else if strings.HasPrefix(query, "prohibited") { //违禁词
		setting.ProhibitedSettingHandler(update, bot.bot)

	} else if strings.Contains(query, "flood_setting") { //刷屏
		setting.FloodSettingHandler(update, bot.bot)

	} else if query == "auto_reply" { //回复消息
		setting.AutoReply(update, bot.bot)

	} else if query == "reply_status_enable" {
		setting.AutoReplyStatus(update, bot.bot, true)

	} else if query == "reply_status_disable" {
		setting.AutoReplyStatus(update, bot.bot, false)

	} else if query == "new_member_check" { //新成员检查
		setting.MemberCheckMenu(update, bot.bot)

	} else if query == "new_member_check_enable" {
		setting.MemberCheckStatus(update, bot.bot, true)

	} else if query == "new_member_check_disable" {
		setting.MemberCheckStatus(update, bot.bot, false)

	} else if query == "new_member_check_time_menu" {
		setting.MemberCheckTimeMenu(update, bot.bot)

	} else if query == "user_check" { //用户检查
		setting.UserCheckMenu(update, bot.bot)

	} else if query == "check_name" {
		setting.NameCheck(update, bot.bot)

	} else if query == "check_username" {
		setting.UserNameCheck(update, bot.bot)

	} else if query == "check_icon" {
		setting.IconCheck(update, bot.bot)

	} else if query == "check_channel" {
		setting.SubScribeCheck(update, bot.bot)

	} else if query == "black_user_list" {
		setting.BlackUserList(update, bot.bot)

	} else if query == "black_user_add" {
		setting.BlackUserAdd(update, bot.bot)

	} else if query == "go_user_check_setting" {
		setting.UserCheckSetting(update, bot.bot)

		//} else if query == "flood_setting" { //刷屏设置
		//	setting.FloodSettingMenu(update, bot.bot)
		//
		//} else if query == "flood_status_enable" {
		//	setting.FloodStatus(update, bot.bot, true)
		//
		//} else if query == "flood_status_disable" {
		//	setting.FloodStatus(update, bot.bot, false)
		//
		//} else if query == "flood_msg_count" {
		//	setting.FloodMsgCountMenu(update, bot.bot)
		//
		//} else if query == "flood_interval" {
		//	setting.FloodIntervalMenu(update, bot.bot)
		//
		//} else if query == "flood_trigger_delete" {
		//	setting.FloodDeleteMsg(update, bot.bot)

	} else if strings.Contains(query, "spam_setting") { //垃圾设置
		setting.SpamSettingHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "dark_model") { //夜晚模式
		setting.DarkSettingHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "punish_setting") { //惩罚设置
		setting.PunishSettingHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "verify_setting") { //入群验证
		group.VerifySettingHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "schedule") { //定时消息
		setting.ScheduleSettingHandler(update, bot.bot)

		// } else if query == "solitaire_home" {
		// 	group.SolitaireHome(update, bot.bot)
	} else if strings.HasPrefix(query, "solitaire_create_step1?") {
		param := query[len("solitaire_create_step1?"):]
		// 创建接龙 step 1
		group.SolitaireCreateStep1(update, bot.bot, param)
	} else if strings.HasPrefix(query, "solitaire_create_limit_time?") {
		param := query[len("solitaire_create_limit_time?"):]
		group.SolitaireCreateStep2LimitTime(update, bot.bot, param)
	} else if strings.HasPrefix(query, "solitaire_create_step2?") {
		// step 2
		param := query[len("solitaire_create_step2?"):]
		group.SolitaireCreateLastStep(update, bot.bot, param)
	} else if strings.HasPrefix(query, "permission_type") { //权限管理
		setting.PermissionSelectHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "permission") { //权限管理
		setting.PermissionHandler(update, bot.bot)

	} else if query == "start" {
		setting.StartHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "manager_group") {
		setting.ManagerGroupHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "delete_notify") {
		setting.DeleteNotifyHandler(update, bot.bot)

	} else {
		msg := tgbotapi.NewMessage(6401399435, "测试推送事件")
		msg.DisableNotification = false
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
	//utils.SendReply(update.CallbackQuery.ID, bot.bot, false, "消息已经处理")
}

func (bot *SmartBot) go_setting(update *tgbotapi.Update) {
	fmt.Println("go_setting...")
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "设置【流量聚集地】群组，选择要更改的项目", utils.SettingMenuMarkup)
	bot.bot.Send(msg)
}
