package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramBot/group"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/setting"
)

// 处理普通消息
func (bot *SmartBot) handleMessage(update *tgbotapi.Update) {
	// 获取用户发送的消息文本

	if group.HandleAdminConversation(update, bot.bot) {
		// 如果消息被拦截 不需要后续处理
		return
	}
	//违禁词处理
	if setting.ProhibitedCheck(update, bot.bot) {
		return
	}

	//规范检查，是否有名字、头像、关联了某个频道
	if setting.UserValidateCheck(update, bot.bot) {
		return
	}

	//是否刷屏
	if setting.FloodCheck(update, bot.bot) {
		return
	}

	//自动回复处理
	if setting.HandlerAutoReply(update, bot.bot) {
		return
	}

	// 保存消息
	message := model.Message{
		ChatId:    update.Message.Chat.ID,
		UserId:    update.Message.From.ID,
		Timestamp: update.Message.Date,
	}
	services.SaveModel(&message, message.ChatId)
}
