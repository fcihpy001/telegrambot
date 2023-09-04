package group

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramBot/model"
	"telegramBot/services"
)

// WelcomeNewMember 进群欢迎
func (mgr *GroupManager) welcomeNewMember(message *tgbotapi.Message) {
	newMembersMsg := message.NewChatMembers
	for _, user := range newMembersMsg {
		if user.IsBot {
			continue
		}
		//读取库里的欢迎语
		welcomeSetting := model.WelcomeSetting{}
		err := services.GetModelData(message.Chat.ID, &welcomeSetting)
		if err != nil {
			logger.Err(err)
			continue
		}
		content := fmt.Sprintf(welcomeSetting.WelcomeText, user.FirstName, message.Chat.Title)

		msg := tgbotapi.NewMessage(message.Chat.ID, content)
		if _, err := mgr.bot.Send(msg); err != nil {
			logger.Err(err)
			continue
		}
	}
}

// CheckUserInfo 检查用户是否是bot, 管理员, 匿名
func (mgr *GroupManager) CheckUserInfo(chatId int64, userId int64) (tgbotapi.ChatMember, error) {
	req := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatId,
			UserID: userId,
		},
	}
	return mgr.bot.GetChatMember(req)
}

// 是否有头像
func (mgr *GroupManager) HasUserProfilePhotos(userId int64) (bool, error) {
	resp, err := mgr.bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: userId,
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		return false, err
	}
	if resp.TotalCount == 0 {
		return false, nil
	}
	// fmt.Println(resp.Photos)
	return true, nil
}

func (mgr *GroupManager) CheckUserIsAdmin(chatId int64, userId int64) (bool, error) {
	info, err := mgr.CheckUserInfo(chatId, userId)
	if err != nil {
		return false, err
	}
	return info.IsAdministrator(), nil
}

func (mgr *GroupManager) fetchAndSaveMember(chatId int64, userId int64) (model.User, error) {
	member, err := mgr.CheckUserInfo(chatId, userId)
	if err != nil {
		return model.User{}, err
	}
	user := model.User{
		Uid:          member.User.ID,
		IsBot:        member.User.IsBot,
		FirstName:    member.User.FirstName,
		LastName:     member.User.LastName,
		Username:     member.User.UserName,
		LanguageCode: member.User.LanguageCode,
	}
	err = services.SaveUser(&user)
	services.UpdateChatMember(chatId, userId, 0, member.Status, 0) // 这里不知道用户是什么时候加入的, 设置为0
	return user, err
}
