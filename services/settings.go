package services

import (
	"log"
	"telegramBot/model"
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
	err := db.Where("chat_id = ?", chatId).First(&setting)
	if err != nil {
		log.Println("get group settings failed")
	}
	return setting
}
