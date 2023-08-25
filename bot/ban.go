package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ban unban

// todo param in text, 4 4min 4day etc
func (bot *SmartBot) Ban(update *tgbotapi.Update) {
	msg := update.Message
	replyTo := msg.ReplyToMessage
	chat := msg.Chat
	if chat == nil {
		log.Printf("invalid message: chat is null: %v", msg)
		return
	}
	chatId := chat.ID
	if replyTo == nil {
		// todo build failed message
		bot.SendText(chatId, "reply message is nil")
		return
	}
	if isAdmin, err := bot.CheckUserIsAdmin(chatId, replyTo.From.ID); isAdmin {
		// todo build failed message
		_ = err
		bot.SendText(chatId, "弄错了吧? 这是管理员哦")
		return
	}

	// if user is administrator
	resp := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID:             chatId,
			SuperGroupUsername: chat.UserName,
			ChannelUsername:    "",
			UserID:             replyTo.From.ID,
		},
		UntilDate:      0,
		RevokeMessages: false,
	}

	bot.sendMessage(resp, "send Ban failed")
}

func (bot *SmartBot) UnBan(update *tgbotapi.Update) {
	msg := update.Message
	replyTo := msg.ReplyToMessage
	chat := msg.Chat
	if chat == nil {
		// todo build failed message
		log.Printf("invalid message: chat is null: %v", msg)
		return
	}
	chatId := chat.ID
	if replyTo == nil {
		// todo build failed message
		bot.SendText(chatId, "reply message is nil")
		return
	}
	if isAdmin, err := bot.CheckUserIsAdmin(chatId, replyTo.From.ID); isAdmin {
		// todo build failed message
		_ = err
		bot.SendText(chatId, "弄错了吧? 这是管理员哦")
		return
	}

	resp := tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID:             chat.ID,
			SuperGroupUsername: chat.UserName,
			UserID:             replyTo.From.ID,
		},
		OnlyIfBanned: true,
	}
	bot.sendMessage(resp, "send UnBan failed")
}

func (bot *SmartBot) SendText(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	bot.bot.Send(msg)
}
