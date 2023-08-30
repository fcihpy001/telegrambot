package services

import (
	"telegramBot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm/clause"
)

func SaveUser(user *model.User) {
	// ignore on duplicate
	err := db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(user).Error
	if err != nil {
		logger.Err(err).Int64("userId", user.Uid).Msg("create user failed")
	}
}

func saveUser(user *tgbotapi.User) {
	SaveUser(&model.User{
		Uid:          user.ID,
		IsBot:        user.IsBot,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.UserName,
		LanguageCode: user.LanguageCode,
	})
}

// 根据 userId 获取 username
func GetUserNames(ids []int64) map[int64]model.User {
	var items []model.User

	users := map[int64]model.User{}
	if err := db.Where("uid in ?", ids).Find(&items).Error; err != nil {
		logger.Err(err).Msg("get user names failed")
		return users
	}

	for _, item := range items {
		users[item.Uid] = item
	}
	return users
}
