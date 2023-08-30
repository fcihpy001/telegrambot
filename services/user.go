package services

import (
	"gorm.io/gorm/clause"
	"telegramBot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
