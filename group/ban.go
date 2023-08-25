package group

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramBot/utils"
)

// ban unban

// todo param in text, 4 4min 4day etc
func (mgr *GroupManager) ban(update *tgbotapi.Update) {
	chatId, fromId, _, err := getChatUserFromReplyMessage(update)
	if err != nil {
		mgr.sendText(chatId, err.Error())
		return
	}

	if isAdmin, err := mgr.CheckUserIsAdmin(chatId, fromId); isAdmin {
		// todo build failed message
		_ = err
		mgr.sendText(chatId, "弄错了吧? 这是管理员哦")
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

	mgr.sendMessage(resp, "send Ban failed")
}

func (mgr *GroupManager) unBan(update *tgbotapi.Update) {
	chatId, fromId, _, err := getChatUserFromReplyMessage(update)
	if err != nil {
		utils.SendText(chatId, err.Error(), mgr.bot)
		return
	}
	if isAdmin, err := mgr.CheckUserIsAdmin(chatId, fromId); isAdmin {
		// todo build failed message
		_ = err
		utils.SendText(chatId, "弄错了吧? 这是管理员哦", mgr.bot)
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
	mgr.sendMessage(resp, "send UnBan failed")
}
