package services

import (
	"log"
	"telegramBot/model"

	"gorm.io/gorm/clause"
)

func SaveSettings(setting *model.WelcomeSetting) {

	if setting.ChatId < 1 {
		return
	}
	//更新或者创建
	if GetSettings(setting.ChatId).ChatId > 0 {
		err := db.Save(setting)
		if err != nil {
			log.Println("update group settings failed", err)
		}
	} else {
		err := db.Create(setting)
		if err != nil {
			log.Println("create group settings failed", err)
		}
	}
}

func GetSettings(chatId int64) model.WelcomeSetting {
	var setting model.WelcomeSetting
	err := db.Where("chat_id = ?", chatId).First(&setting).Error
	if err != nil {
		logger.Err(err).Msg("get group settings failed")
	}
	return setting
}

func SaveInviteSettings(setting *model.InviteSetting) {
	if setting.ChatId < 1 {
		logger.Error().Int64("chatId", setting.ChatId).Msg("invalild chatId")
		return
	}
	//更新或者创建
	if GetInviteSettings(setting.ChatId).ChatId > 0 {
		err := db.Save(setting)
		if err != nil {
			log.Println("update invite settings failed", err)
		}
	} else {
		err := db.Create(setting)
		if err != nil {
			log.Println("create invite settings failed", err)
		}
	}
}

func GetInviteSettings(chatId int64) model.InviteSetting {
	var setting model.InviteSetting
	err := db.Where("chat_id = ?", chatId).First(&setting)
	if err != nil {
		log.Println("get invite settings failed")
	}
	return setting
}

func SaveReplySettings(model *model.ReplySetting) {
	if model.ChatId < 1 {
		return
	}
	//更新或者创建
	if GetReplySettings(model.ChatId).ChatId > 0 {
		err := db.Save(model)
		if err != nil {
			log.Println("update reply settings failed", err)
		}
	} else {
		err := db.Create(model)
		if err != nil {
			log.Println("create reply settings failed", err)
		}
	}
}

func GetReplySettings(chatId int64) model.ReplySetting {
	var setting model.ReplySetting
	err := db.Where("chat_id = ?", chatId).First(&setting)
	if err != nil {
		log.Println("get reply settings failed")
	}
	return setting
}

func SaveProhibitSettings(model *model.ProhibitedSetting) {
	if model.ChatId < 1 {
		return
	}
	//更新或者创建
	if GetProhibitSettings(model.ChatId).ChatId > 0 {
		err := db.Save(model)
		if err != nil {
			log.Println("update Prohibit settings failed", err)
		}
	} else {
		err := db.Create(model)
		if err != nil {
			log.Println("create Prohibit settings failed", err)
		}
	}
}

func GetProhibitSettings(chatId int64) model.ProhibitedSetting {
	var setting model.ProhibitedSetting
	err := db.Where("chat_id = ?", chatId).First(&setting).Error
	if err != nil {
		log.Println("get Prohibit settings failed")
	}
	return setting
}

func SaveMemberSettings(model *model.NewMemberCheck) {
	if model.ChatId < 1 {
		return
	}
	//更新或者创建
	if GetMemberSettings(model.ChatId).ChatId > 0 {
		err := db.Save(model).Error
		if err != nil {
			log.Println("update Prohibit settings failed", err)
		}
	} else {
		err := db.Create(model).Error
		if err != nil {
			log.Println("create Prohibit settings failed", err)
		}
	}
}

func GetMemberSettings(chatId int64) model.NewMemberCheck {
	var setting model.NewMemberCheck
	err := db.Where("chat_id = ?", chatId).First(&setting).Error
	if err != nil {
		log.Println("get Prohibit settings failed")
	}
	return setting
}

func SaveModel(model interface{}, chatId int64) {
	if chatId == 0 {
		return
	}
	err := db.Save(model)
	if err != nil {
		log.Println("update or insert model data failed", err)
	}
}

func GetModelData(chatId int64, model interface{}) error {

	err := db.Where("chat_id = ?", chatId).First(&model)
	if err != nil {
		log.Println("get Prohibit settings failed")
		return err.Error
	}
	return nil
}

func GetAllProhibitSettings() ([]model.ProhibitedSetting, error) {
	var items []model.ProhibitedSetting
	err := db.Find(&items).Error
	return items, err
}

func GetAllUserCheck() ([]model.UserCheck, error) {
	var items []model.UserCheck
	err := db.Find(&items).Error
	return items, err
}

func GetAllCautions() ([]model.UserCautions, error) {
	var items []model.UserCautions

	err := db.Find(&items).Error
	return items, err
}

func UpdateUserCaution(chatId, userId int64, triggerType model.TriggerType, count int) error {
	item := model.UserCautions{
		ChatId:       chatId,
		UserId:       userId,
		TriggerType:  string(triggerType),
		TriggerCount: int64(count),
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "chat_id"}, {Name: "trigger_type"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"trigger_count": count,
		}),
	}).Create(&item).Error
}
