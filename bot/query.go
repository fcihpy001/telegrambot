package bot

import (
	"fmt"
	"log"
	"strings"
	"telegramBot/group"
	"telegramBot/lucky"
	"telegramBot/setting"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 处理行内按钮事件
func (bot *SmartBot) handleQuery(update *tgbotapi.Update) {
	query := update.CallbackQuery.Data
	fmt.Println("query command--", query)

	if strings.HasPrefix(query, "lucky") {
		lucky.LuckyHandler(update, bot.bot)

	} else if query == "start" { //开始界面
		setting.StartHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "settings") { //设置
		setting.Settings(update, bot.bot)

	} else if query == "go_setting" { //设置主菜单
		bot.go_setting(update)

	} else if strings.HasPrefix(query, "group") { //群组设置
		group.GroupHandlerQuery(update, bot.bot)

	} else if strings.HasPrefix(query, "prohibited") { //违禁词
		setting.ProhibitedSettingHandler(update, bot.bot)

	} else if strings.Contains(query, "flood_setting") { //刷屏
		setting.FloodSettingHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "auto_reply") { //回复消息
		setting.ReplyHandler(update, bot.bot)

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

	} else if strings.HasPrefix(query, "permission_type") { //权限管理
		setting.PermissionSelectHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "permission") { //权限管理
		setting.PermissionHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "manager_group") { //群组管理
		setting.ManagerGroupHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "delete_notify") { //删除通知
		setting.DeleteNotifyHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "user_check") { //用户检查
		setting.UserCheckHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "new_member_check") { //新成员检查
		setting.MemberCheckHandler(update, bot.bot)

	} else if query == "solitaire_home" {
		group.SolitaireHome(update, bot.bot)
	} else if strings.HasPrefix(query, "solitaire_create_step1?") {
		param := query[len("solitaire_create_step1?"):]
		// 创建接龙 step 1
		group.SolitaireCreateStep1(update, bot.bot, param)
	} else if strings.HasPrefix(query, "solitaire_create_step2?") {
		param := query[len("solitaire_create_step2?"):]
		group.SolitaireCreateStep2(update, bot.bot, param)
	} else if strings.HasPrefix(query, "solitaire_create_limit_time?") {
		param := query[len("solitaire_create_limit_time?"):]
		group.SolitaireCreateStep2LimitTime(update, bot.bot, param)
	} else if strings.HasPrefix(query, "solitaire_create_last_step?") {
		// last step
		param := query[len("solitaire_create_last_step?"):]
		group.SolitaireCreateLastStep(update, bot.bot, param)
	} else if strings.HasPrefix(query, "solitaire_delete?") {
		// delete solitaire
		param := query[len("solitaire_delete?"):]
		group.SolitaireDelete(update, bot.bot, param)
	} else if strings.HasPrefix(query, "solitaire_confirm_delete?") {
		param := query[len("solitaire_confirm_delete?"):]
		group.SolitaireConfirmDelete(update, bot.bot, param)

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
