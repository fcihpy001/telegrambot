package services

import (
	"fmt"
	"log"
	"telegramBot/model"
	"telegramBot/utils"

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

func GetProhibitSettings(chatId int64) model.ProhibitedSetting {
	var setting model.ProhibitedSetting
	err := db.Model(&model.ProhibitedSetting{}).Preload("Punishment").Where("chat_id = ?", chatId).First(&setting).Error
	if err != nil {
		log.Println("get Prohibit settings failed")
	}
	return setting
}

func SaveModel(model interface{}, chatId int64) {
	if chatId == 0 {
		fmt.Println("不符合存储条件")
		return
	}
	err := db.Save(model)
	if err.Error != nil {
		log.Println("update or insert model data failed", err)
	}
}

func SaveModelWhere(model interface{}, chatId int64, where string) {
	if chatId == 0 {
		fmt.Println("不符合存储条件")
		return
	}
	err := db.Where(where).Save(model)
	if err.Error != nil {
		log.Println("update or insert model data failed", err)
	}
}

func GetModelData(chatId int64, model interface{}) error {

	err := db.Where("chat_id = ?", chatId).First(&model)
	if err.Error != nil {
		log.Println("get model data  failed", err.Error)
		return err.Error
	}
	fmt.Println("get model data success::", model)
	return nil
}
func GetModels(models []interface{}) error {
	result := db.Find(&models)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAllDarkSettings() ([]model.DarkModelSetting, error) {
	var items []model.DarkModelSetting
	err := db.Find(&items).Error
	return items, err
}

func GetReplySetting(chat_id int64) (model.ReplySetting, error) {
	var items model.ReplySetting
	err := db.Model(&model.ReplySetting{}).Preload("ReplyList").Where("chat_id = ?", chat_id).First(&items).Error
	return items, err
}

func GetAllReply(chat_id int64) ([]model.Reply, error) {
	var items []model.Reply
	err := db.Where("chat_id = ?", chat_id).Find(&items).Error
	return items, err
}

func DeleteReply(keyword string, chat_id int64) error {
	var item model.Reply
	err := db.Where("chat_id = ? and key_world = ?", chat_id, keyword).Delete(&item).Error
	return err
}

func GetAllTask(chat_id int64) ([]model.Task, error) {
	var items []model.Task
	err := db.Find(&items).Error
	return items, err
}

func GetAllDeleteTask() ([]model.ScheduleDelete, error) {
	var items []model.ScheduleDelete
	err := db.Find(&items).Error
	return items, err
}

func DeleteTask(task *model.ScheduleDelete) error {
	err := db.Delete(&task).Error
	return err
}

func DeleteModel(model interface{}, where string) error {
	err := db.Where(where).Delete(&model).Error
	return err
}

func DeleteInviteData() error {
	var item model.InviteSetting
	err := db.Where("chat_id =?", utils.GroupInfo.GroupId).Delete(&item).Error
	return err
}

func GetAllGroups(where string) ([]model.GroupInfo, error) {
	var items []model.GroupInfo
	err := db.Where(where).Find(&items).Error
	return items, err
}

func GetModelWhere(where string, model interface{}) error {

	err := db.Where(where).First(&model)
	if err.Error != nil {
		log.Println("get model failed")
		return err.Error
	}
	return nil
}

func GetScheduleMsgList(where string) ([]model.ScheduleMsg, error) {
	var items []model.ScheduleMsg
	err := db.Where(where).Find(&items).Error
	if len(where) == 0 {
		err = db.Find(&items).Error
	}
	return items, err
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

func Update2GroupInfoCaution(chatId, userId int64, triggerType model.TriggerType, count int) error {
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
