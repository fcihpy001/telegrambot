package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// WelcomeNewMember 进群欢迎
func (bot *SmartBot) WelcomeNewMember(message *tgbotapi.Message) {
	newMembersMsg := message.NewChatMembers

	for _, user := range newMembersMsg {
		if user.IsBot {
			continue
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, "👏👏👏 欢迎 "+user.FirstName+" 加入"+message.Chat.Title)
		if _, err := bot.bot.Send(msg); err != nil {
			log.Println(err)
			continue
		}
	}
}

// CheckUserInfo 检查用户是否是bot, 管理员, 匿名
func (bot *SmartBot) CheckUserInfo(chatId int64, userId int64) (tgbotapi.ChatMember, error) {
	req := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatId,
			UserID: userId,
		},
	}
	return bot.bot.GetChatMember(req)
}

func (bot *SmartBot) CheckUserIsAdmin(chatId int64, userId int64) (bool, error) {
	info, err := bot.CheckUserInfo(chatId, userId)
	if err != nil {
		return false, err
	}
	return info.IsAdministrator(), nil
}
