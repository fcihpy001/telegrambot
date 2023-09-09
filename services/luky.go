package services

import (
	"telegramBot/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateLucky(item *model.LuckyActivity) {
	if err := db.Save(item).Error; err != nil {
		logger.Err(err).Msg("create lucky activity failed")
	}
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

func GetAllLuckyActivities(chatId int64, offset, limit int) (items []*model.LuckyActivity) {
	err := db.Where("chat_id = ?", chatId).Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		log.Err(err).Msg("find lucky activity failed")
	}
	return
}

func GetLuckyActivity(chatId int64, offset int) (item *model.LuckyActivity) {
	var items []*model.LuckyActivity
	err := db.Where("chat_id = ?", chatId).Offset(offset).Limit(1).Find(&items).Error
	if err != nil {
		log.Err(err).Msg("find lucky activity failed")
		return
	}
	if len(items) > 0 {
		item = items[0]
	}
	return
}

func GetLuckyActivityCount(chatId int64) int64 {
	var count int64
	err := db.Model(&model.LuckyActivity{}).Where("chat_id = ?", chatId).Count(&count).Error
	if err != nil {
		logger.Err(err).Msg("count lucky activity failed")
	}
	return count
}
