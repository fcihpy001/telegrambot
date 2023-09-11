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

// 获取所有正在进行的抽奖
func GetAllLuckyActivities() (items []*model.LuckyActivity) {
	err = db.Where("status = ?", model.LuckyStatusStart).Find(&items).Error
	return
}

func GetChatGroupLuckyActivities(chatId int64, offset, limit int) (items []*model.LuckyActivity) {
	var err error
	if limit > 0 {
		err = db.Where("chat_id = ?", chatId).Offset(offset).Limit(limit).Find(&items).Error
	} else {
		err = db.Where("chat_id = ?", chatId).Offset(offset).Find(&items).Error
	}

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

func UpdateLuckyActivity(record *model.LuckyActivity) {
	if err := db.Save(record).Error; err != nil {
		logger.Err(err).Msg("update lucky activity failed")
	}
}

func GetLuckyActivityCount(chatId int64) int64 {
	var count int64
	err := db.Model(&model.LuckyActivity{}).Where("chat_id = ?", chatId).Count(&count).Error
	if err != nil {
		logger.Err(err).Msg("count lucky activity failed")
	}
	return count
}

func StatChatLuckyCount(chatId int64) (total, opened, canceled int) {
	var items []*model.LuckyActivity
	err := db.Model(&model.LuckyActivity{}).Where("chat_id = ?", chatId).Find(&items).Error
	if err != nil {
		return 0, 0, 0
	}

	total = len(items)
	for _, item := range items {
		if item.Status == model.LuckyStatusEnd {
			opened++
		} else if item.Status == model.LuckyStatusCancel {
			canceled++
		}
	}
	return
}

func OnLuckyParticipate(record *model.LuckyActivity, fromId int64, username string) {
	item := model.LuckyRecord{
		LuckyId:  int64(record.ID),
		ChatId:   record.ChatId,
		UserId:   fromId,
		Username: username,
	}
	if err := db.Save(&item).Error; err != nil {
		logger.Err(err).Msg("save lucky participate record failed")
	}
	if err := db.Save(record).Error; err != nil {
		logger.Err(err).Msg("update lucky activity failed")
	}
}

// 用户是否已经参与过
func CheckUserHasParticipated(luckyId, userId int64) bool {
	var count int64
	err := db.Model(&model.LuckyRecord{}).Where("lucky_id = ? AND user_id = ?", luckyId, userId).Count(&count).Error
	if err != nil {
		logger.Err(err).Msg("CheckUserHasParticipated failed")
	}
	return count > 0
}

func GetLuckyAllParticipates(record *model.LuckyActivity) []model.LuckyRecord {
	var parts []model.LuckyRecord

	if err := db.Where("lucky_id = ?", record.ID).Find(&parts).Error; err != nil {
		logger.Err(err).Msg("gete participates failed")
	}

	return parts
}

func UpdateLuckyRewardRecord(record *model.LuckyRecord) {
	if err := db.Save(record).Error; err != nil {
		logger.Err(err).Msg("update lucky reward failed")
	}
}
