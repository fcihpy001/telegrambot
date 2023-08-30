package services

import (
	"fmt"
	"telegramBot/model"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm/clause"
)

func saveChatGroup(chat *tgbotapi.Chat) {
	photo := ""
	location := ""
	if chat.Photo != nil {
		photo = chat.Photo.SmallFileID
	}
	if chat.Location != nil {
		location = chat.Location.Address
	}
	SaveChatGroup(chat.ID, chat.Title, chat.Type, chat.UserName, photo, location)
}

func SaveChatGroup(chatId int64, title, typ, username, photo, location string) {
	chat := model.ChatGroup{
		ChatId:    chatId,
		Title:     title,
		GroupType: typ,
		Chatname:  username,
		Photo:     photo,
		Location:  location,
	}
	// ignore on duplicate
	err := db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&chat).Error
	if err != nil {
		logger.Err(err).Int64("chatId", chatId).Msg("create chat group failed")
	}
}

// create or update chat member
func UpdateChatMember(chatId, userId int64, status string) {
	item := model.UserChat{
		UserId: userId,
		ChatId: chatId,
		Status: status,
	}
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "chat_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"stats": status}),
	}).Create(&item).Error
	if err != nil {
		logger.Err(err).
			Int64("chatId", chatId).
			Int64("userId", userId).
			Str("status", status).
			Msg("update user-chat member status failed")
	}
}

func RemoveChatMember(chatId, userId int64) {
	db.Exec("delete user_chat where user_id = ? and chat_id = ?", userId, chatId)
}

func SaveUserAction(userId, chatId int64, action string) {
	now := time.Now()
	db.Save(&model.UserAction{
		Action: action,
		UserId: userId,
		ChatId: chatId,
		Day:    fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day()),
	})
}
