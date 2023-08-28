package services

import "telegramBot/model"

func SaveUser(user *model.User) {
	db.Create(user)
}
