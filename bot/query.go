package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/group"
	"telegramBot/lucky"
	"telegramBot/model"
	"telegramBot/setting"
	"telegramBot/utils"
)

// 处理行内按钮事件
func (bot *SmartBot) handleQuery(update *tgbotapi.Update) {
	query := update.CallbackQuery.Data
	fmt.Println("query command--", query)

	if strings.HasPrefix(query, "lucky") {
		lucky.LuckyHandler(update, bot.bot)

	} else if strings.HasPrefix(query, "group") {
		group.GroupHandlerQuery(update, bot.bot)

	} else if strings.HasPrefix(query, "settings") {
		setting.Settings(update.CallbackQuery.Message.Chat.ID, bot.bot)

	} else if query == "join_group" {
		fmt.Println("replay...")
		// 创建 ForceReply 结构
		forceReply := tgbotapi.ForceReply{
			ForceReply: true,
		}

		// 创建包含 ForceReply 的消息
		message := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "请回复此消息：")
		message.ReplyMarkup = forceReply

		// 发送消息
		_, err := bot.bot.Send(message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	} else if query == "next_page" {
		//	发送还键盘的推送消息
		msg := tgbotapi.NewMessage(6401399435, "🎁【零度社区 (LingduDAO)- 中文群】群组发起了发言次数抽奖活动\n已开奖：1       未开奖：1       取消：0\n\nLDD是零度DAO的社区币\n├参与条件：发言6条\n├发言起始统计时间：2023-08-28 11:20:00\n├开奖时间：2023-08-28 22:00:00\n├奖品列表：\n├       2USDT     ×3份\n\n【如何参与？】在群组中发言6次，参与活动。")
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🌺加入活动群众", "lucky_activity"),
			))
		msg.ReplyMarkup = inlineKeyboard
		msg.DisableNotification = false
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if query == "go_setting" {
		bot.go_setting(update)

	} else if query == "prohibited_words" { //违禁词
		setting.ProhibitedSettingHandler(update, bot.bot)

	} else if query == "prohibited_list" { //违禁词列表
		setting.ProhibitedList(update, bot.bot)
	} else if query == "prohibited_add_menu" {
		setting.ProhibitedAddMenu(update, bot.bot)

	} else if query == "prohibited_delete_menu" {
		setting.ProhibitedDeleteMenu(update, bot.bot)

	} else if query == "prohibited_delete" { //违禁词开关
		setting.ProhibitedDelete(update, bot.bot)

	} else if query == "prohibitedStatus_enable" {
		setting.ProhibitedStatus(update, bot.bot, true)

	} else if query == "prohibitedStatus_disable" {
		setting.ProhibitedStatus(update, bot.bot, false)

	} else if query == "prohibited_punish_setting" {
		setting.PunishSetting(update, bot.bot)

	} else if query == "prohibit_punish_type1" {
		setting.PunishAction(update, bot.bot, model.PunishTypeWarning)

	} else if query == "prohibit_punish_type2" {
		setting.PunishAction(update, bot.bot, model.PunishTypeBan)

	} else if query == "prohibit_punish_type3" {
		setting.PunishAction(update, bot.bot, model.PunishTypeKick)

	} else if query == "prohibit_punish_type4" {
		setting.PunishAction(update, bot.bot, model.PunishTypeBanAndKick)

	} else if query == "prohibit_punish_type5" {
		setting.PunishAction(update, bot.bot, model.PunishTypeRevoke)

	} else if query == "prohibit_warning_count1" {
		setting.WarningCount(update, bot.bot, 1)

	} else if query == "prohibit_warning_count2" {
		setting.WarningCount(update, bot.bot, 2)

	} else if query == "prohibit_warning_count3" {
		setting.WarningCount(update, bot.bot, 3)

	} else if query == "prohibit_warning_count4" {
		setting.WarningCount(update, bot.bot, 4)

	} else if query == "prohibit_warning_count5" {
		setting.WarningCount(update, bot.bot, 5)

	} else if query == "prohibit_warning_after_action1" {
		setting.WarningAction(update, bot.bot, model.PunishTypeBan)

	} else if query == "prohibit_warning_after_action2" {
		setting.WarningAction(update, bot.bot, model.PunishTypeKick)

	} else if query == "prohibit_warning_after_action3" {
		setting.WarningAction(update, bot.bot, model.PunishTypeBanAndKick)

	} else if query == "go_prohibited_setting" {
		setting.GoProhibitedSetting(update, bot.bot)

	} else if query == "prohibited_ban_time" {
		setting.PunishTime(update, bot.bot)

	} else if query == "prohibited_ban_time_type1" {
		setting.PunishTimeType(update, bot.bot, model.BanTimeType1)
	} else if query == "prohibited_ban_time_type2" {
		setting.PunishTimeType(update, bot.bot, model.BanTimeType2)
	} else if query == "prohibited_ban_time_type3" {
		setting.PunishTimeType(update, bot.bot, model.BanTimeType3)
	} else if query == "prohibited_ban_time_type4" {
		setting.PunishTimeType(update, bot.bot, model.BanTimeType4)
	} else if query == "prohibited_ban_time_type5" {
		setting.PunishTimeType(update, bot.bot, model.BanTimeType5)
	} else if query == "prohibited_ban_time_type6" {
		setting.PunishTimeType(update, bot.bot, model.BanTimeType6)

	} else if query == "auto_reply" { //回复消息
		setting.AutoReply(update, bot.bot)

	} else if query == "reply_status_enable" {
		setting.AutoReplyStatus(update, bot.bot, true)

	} else if query == "reply_status_disable" {
		setting.AutoReplyStatus(update, bot.bot, false)

	} else if query == "new_member_check" {
		setting.MemberCheckMenu(update, bot.bot)

	} else if query == "new_member_check_enable" {
		setting.MemberCheckStatus(update, bot.bot, true)

	} else if query == "new_member_check_disable" {
		setting.MemberCheckStatus(update, bot.bot, false)

	} else if query == "new_member_check_time_menu" {
		setting.MemberCheckTimeMenu(update, bot.bot)

	} else if query == "user_check" {
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
