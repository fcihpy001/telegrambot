package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Date when restrictions will be lifted for the user, unix time.
// If user is restricted for more than 366 days or less than 30 seconds
// from the current time, they are considered to be restricted forever
func (bot *SmartBot) Mute(update *tgbotapi.Update) {
	chatId, fromId, _, err := getChatUserFromReplyMessage(update)
	if err != nil {
		bot.SendText(chatId, err.Error())
		return
	}

	if isAdmin, _ := bot.CheckUserIsAdmin(chatId, fromId); isAdmin {
		bot.SendText(chatId, "弄错了吧? 这是管理员哦")
		return
	}
	// parse until param from message text
	seconds := parseUntilDate(update.Message.Text)
	msg := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: fromId,
		},
		UntilDate: seconds,
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       false,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
		},
	}
	bot.sendMessage(msg, "mute user failed")
}

func (bot *SmartBot) UnMute(update *tgbotapi.Update) {
	chatId, fromId, _, err := getChatUserFromReplyMessage(update)
	if err != nil {
		bot.SendText(chatId, err.Error())
		return
	}

	if isAdmin, _ := bot.CheckUserIsAdmin(chatId, fromId); isAdmin {
		bot.SendText(chatId, "弄错了吧? 这是管理员哦")
		return
	}
	msg := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: fromId,
		},
		UntilDate: 0, // todo
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
			CanChangeInfo:         true,
			CanInviteUsers:        true,
		},
	}
	bot.sendMessage(msg, "unmute user failed")
}
