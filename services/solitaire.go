package services

import (
	"context"
	"telegramBot/model"
	"time"

	"gorm.io/gorm"
)

// 接龙

// 获取群组接龙列表和
func GetChatSolitaireList(chatId int64) ([]model.Solitaire, error) {
	var items []model.Solitaire
	err := db.Where("chat_id = ? AND status = ?", chatId, model.SolitaireStatusActive).Find(&items).Error
	return items, err
}

func GetChatSolitaireById(sid int) (item model.Solitaire, err error) {
	err = db.First(&item, sid).Error
	return
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

func GetSolitaireMessageList(sid int64) (items []model.SolitaireMessage, err error) {
	if err = db.Where("solitaire_id = ?", sid).Find(&items).Error; err != nil {
		logger.Err(err).Msg("find solitarie message list failed")
	}
	return
}

// 每个接龙每个用户只能参与一次
func GetUserSolitaire(sid int, userId int64) (*model.SolitaireMessage, error) {
	var msg model.SolitaireMessage

	err := db.Where("solitaire_id = ? AND user_id = ?", sid, userId).First(&msg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

func UpdateSolitaireStatusByIncCollected(sid uint) error {
	err := db.Exec("UPDATE solitaire s SET s.msg_collected = s.msg_collected + 1, status = (case when s.msg_collected >= s.limit_users and s.limit_users > 0  then 'end' else s.status end) where id = ? ", sid).Error
	if err != nil {
		logger.Err(err).Uint("solitaireId", sid).Msg("update solitaire status failed")
	}
	return err
}

// 每分钟检查所有接龙 是否超时
func CheckSolitaireEnded(ctx context.Context) {
	tkr := time.NewTicker(time.Minute)
	var items []model.Solitaire

	for {
		select {
		case <-tkr.C:
			err := db.Where("limit_time > 0 AND status = ? AND deleted_at IS NOT NULL", "active").Find(&items).Error
			if err == nil {
				now := time.Now().Unix()
				for _, item := range items {
					if int64(item.LimitTime) < now {
						// update status
						item.Status = "end"
						db.Save(&item)
					}
				}
			} else {
				logger.Err(err).Msg("CheckSolitaireEnded: get solitaire list failed")
			}

		case <-ctx.Done():
			logger.Info().Msg("CheckSolitaireEnded routine exit")
			return
		}
	}
}
