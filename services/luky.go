package services

import (
	"telegramBot/model"

	"gorm.io/gorm"
)

func CreateLucky() {

}

func FindChatLuckySetting(chatId int64) *model.LuckySetting {
	var setting model.LuckySetting
	err := db.Where("chat_id = ? AND deleted_at IS NULL", chatId).First(&setting).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Err(err).Int64("chatId", chatId).Msg("FindChatLuckySetting failed")
		}
		return nil
	}
	return &setting
}

func UpdateChatLuckySettings(item *model.LuckySetting) {
	if err := db.Save(item).Error; err != nil {
		logger.Err(err).Msg("update or save lucky setting failed")
	}
}
