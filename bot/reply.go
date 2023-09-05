package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegramBot/group"
	"telegramBot/setting"
)

// 处理需要用户回复的消息，如请输入名字。。。等
func (bot *SmartBot) handleReply(update *tgbotapi.Update) {

	replyMsg := update.Message.ReplyToMessage.Text
	if strings.Contains(replyMsg, "输入添加的违禁词（一行一个") {
		setting.ProhibitedAddResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "请输入新群员限制时间") {
		setting.MemberCheckTimeResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "请输入要禁止的名字（一行一个") {
		setting.BlackUserAddResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "请输入时间内发送消息的最大条数") {
		setting.FloodMsgCountResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "请输入统计发送消息的间隔时间（秒）") {
		setting.FloodIntervalResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "输入允许的姓名最大长度（例如：15") {
		setting.SpamNameLengthReply(update, bot.bot)

	} else if strings.Contains(replyMsg, "输入允许的消息最大长度") {
		setting.SpamMsgLengthReply(update, bot.bot)

	} else if strings.Contains(replyMsg, "输入你想要设置内容：") {
		setting.ScheduleAndContentResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "在开启状态下，到达设定时间才会发送消息，请回复开始时间") {
		setting.ScheduleDateStartResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "到达设定时间后自动停止，请回复终止时间") {
		setting.ScheduleDateEndResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "输入要设置的新成员入群欢迎内容，占位符中%s代替") {
		group.WelcomeTextSettingResult(update, bot.bot)

	} else if strings.Contains(replyMsg, "👉 输入处罚禁言的时长（分钟，例如：60") {
		setting.BanTimeReply(update, bot.bot)

	} else if strings.Contains(replyMsg, "请输入要删除的违禁词（一行一个）") {
		setting.ProhibitedDeleteResult(update, bot.bot)
	}
}
