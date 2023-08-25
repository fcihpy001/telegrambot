package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ban unban

// todo param in text, 4 4min 4day etc
func (bot *SmartBot) Ban(update *tgbotapi.Update) {
	chatId, fromId, _, err := getChatUserFromReplyMessage(update)
	if err != nil {
		bot.SendText(chatId, err.Error())
		return
	}

	if isAdmin, err := bot.CheckUserIsAdmin(chatId, fromId); isAdmin {
		// todo build failed message
		_ = err
		bot.SendText(chatId, "弄错了吧? 这是管理员哦")
		return
	}
	seconds := parseUntilDate(update.Message.Text)
	// if user is administrator
	resp := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			// SuperGroupUsername: chat.UserName,
			ChannelUsername: "",
			UserID:          fromId,
		},
		UntilDate:      seconds,
		RevokeMessages: false,
	}

	bot.sendMessage(resp, "send Ban failed")
}

func (bot *SmartBot) UnBan(update *tgbotapi.Update) {
	chatId, fromId, _, err := getChatUserFromReplyMessage(update)
	if err != nil {
		bot.SendText(chatId, err.Error())
		return
	}
	if isAdmin, err := bot.CheckUserIsAdmin(chatId, fromId); isAdmin {
		// todo build failed message
		_ = err
		bot.SendText(chatId, "弄错了吧? 这是管理员哦")
		return
	}

	resp := tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			// SuperGroupUsername: chat.UserName,
			UserID: fromId,
		},
		OnlyIfBanned: true,
	}
	bot.sendMessage(resp, "send UnBan failed")
}

func (bot *SmartBot) SendText(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	bot.bot.Send(msg)
}
