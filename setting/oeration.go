package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegramBot/utils"
	"time"
)

func OperationHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	cmd := update.Message.Command()
	switch cmd {
	case "kick":
		kickUserHandler(update, bot)

	case "ban":
		banUserHandler(update, bot)

	case "unban":
		unBanUserHandler(update, bot)

	case "mute":
		muteUserHandler(update, bot, 10*60)

	case "unmute":
		muteUserHandler(update, bot, 0)

	}
}

// 对用户进行禁言,second=解除禁言
func muteUserHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, second int) {

	chatId := update.Message.Chat.ID
	userId := update.Message.ReplyToMessage.From.ID
	if update.Message.ReplyToMessage == nil {
		utils.SendText(chatId, "请在要操作的用户所发的消息上，回复此命令", bot)
		return
	}
	MuteUser(chatId, bot, second, userId)
}

func MuteUser(chatId int64, bot *tgbotapi.BotAPI, second int, userId int64) {

	permission := &tgbotapi.ChatPermissions{
		CanSendMessages:       true,
		CanSendMediaMessages:  true,
		CanSendPolls:          true,
		CanSendOtherMessages:  true,
		CanAddWebPagePreviews: true,
		CanChangeInfo:         true,
		CanInviteUsers:        true,
		CanPinMessages:        true,
	}
	if second > 0 {
		permission = &tgbotapi.ChatPermissions{
			CanSendMessages:       false,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
			CanPinMessages:        false,
		}
	}
	date := time.Now().Add(time.Duration(second) * time.Second).Unix()
	msg := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
		UntilDate:   date,
		Permissions: permission,
	}
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func banMember(bot *tgbotapi.BotAPI, chatId int64, second int, userId int64, muteMedia bool) {
	permission := &tgbotapi.ChatPermissions{
		CanSendMessages:       false,
		CanSendMediaMessages:  false,
		CanSendPolls:          false,
		CanSendOtherMessages:  false,
		CanAddWebPagePreviews: false,
		CanChangeInfo:         false,
		CanInviteUsers:        false,
		CanPinMessages:        false,
	}
	if muteMedia {
		permission = &tgbotapi.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
			CanPinMessages:        false,
		}
	}
	date := time.Now().Add(time.Duration(second) * time.Second).Unix()
	msg := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
		UntilDate:   date,
		Permissions: permission,
	}
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// Pin某条消息
func PinMessage(chatId int64, bot *tgbotapi.BotAPI, messageId int) {
	msg := tgbotapi.PinChatMessageConfig{
		ChatID:              chatId,
		ChannelUsername:     "",
		MessageID:           messageId,
		DisableNotification: false,
	}
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func UnPinMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI, messageId int) {
	msg := tgbotapi.UnpinChatMessageConfig{
		ChatID:          update.Message.Chat.ID,
		ChannelUsername: "",
		MessageID:       messageId,
	}
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// 踢出某个用户
func kickUserHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatId := update.Message.Chat.ID
	if update.Message.ReplyToMessage == nil {
		utils.SendText(chatId, "请在要操作的用户所发的消息上，回复此命令", bot)
		return
	}
	userId := update.Message.ReplyToMessage.From.ID
	kickUser(update, bot, userId)
}

func kickUser(update *tgbotapi.Update, bot *tgbotapi.BotAPI, userId int64) {

	msg := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: update.Message.Chat.ID,
			UserID: userId,
		},
		UntilDate:      time.Now().Add(7 * 60 * 24 * time.Minute).Unix(),
		RevokeMessages: false,
	}
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func banUserHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatId := update.Message.Chat.ID
	if update.Message.ReplyToMessage == nil {
		utils.SendText(chatId, "请在要操作的用户所发的消息上，回复此命令", bot)
		return
	}
	userId := update.Message.ReplyToMessage.From.ID
	banUser(update, bot, userId)
}

func banUser(update *tgbotapi.Update, bot *tgbotapi.BotAPI, userId int64) {

	msg := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: update.Message.Chat.ID,
			UserID: userId,
		},
		UntilDate:      time.Now().Add(1 * 60 * 24 * time.Minute).Unix(),
		RevokeMessages: false,
	}
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func unBanUserHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatId := update.Message.Chat.ID
	if update.Message.ReplyToMessage == nil {
		utils.SendText(chatId, "请在要操作的用户所发的消息上，回复此命令", bot)
		return
	}
	msg := tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: update.Message.ReplyToMessage.From.ID,
		},
		OnlyIfBanned: false,
	}
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
