package services

import (
	"telegramBot/model"

	"gorm.io/gorm"
)

// 接龙

// 获取群组接龙列表和
func GetChatSolitaireList(chatId int64) ([]model.Solitaire, error) {
	var items []model.Solitaire
	err := db.Where("chat_id = ? AND status = ?", chatId, model.SolitaireStatusActive).Find(&items).Error
	return items, err
}

func CreateSolitaire(chatId, creator int64, limitUsers int, limitTime int64, description string) (model.Solitaire, error) {
	item := model.Solitaire{
		ChatId:       chatId,
		LimitUsers:   limitUsers,
		LimitTime:    int(limitTime), // 截止时间
		Creator:      creator,        // 创建用户 id
		Description:  description,    // `gorm:""`
		MsgCollected: 0,              // 已接龙条数
		Status:       model.SolitaireStatusActive,
	}
	err := db.Save(&item).Error
	return item, err
}

func DeleteSolitaire(sid int64) {
	db.Delete(&model.Solitaire{}, sid)
}

func NewChatSolitaireMessage(chatId, sid, userId int64, msg string) {
	// 1. update solitarie
	var item model.Solitaire
	item.ID = uint(sid)
	if err := db.Model(&item).Update("msg_collected", gorm.Expr("msg_collected + 1")).Error; err != nil {
		logger.Err(err).Int64("solitaireId", sid).Int64("chatId", chatId).Msg("update solitaire collected failed")
	}
	// 2. insert solitarie detail
	err := db.Save(&model.SolitaireMessage{
		ChatId:      chatId,
		SolitaireId: sid,
		UserId:      userId,
		Message:     msg,
	}).Error
	if err != nil {
		logger.Err(err).Msg("create solitaire failed")
	}
}

func GetSolitaireMessageList(sid int64) (items []model.SolitaireMessage) {
	if err := db.Where("solitaire_id = ?", sid).Find(&items).Error; err != nil {
		logger.Err(err).Msg("find solitarie message list failed")
	}
	return
}
