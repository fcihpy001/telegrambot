package services

import (
	"context"
	"fmt"
	"telegramBot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

/**
把每分钟各个群的各个用户的消息数量统计到redis中, 定时将上一分钟的统计数据写入db并清除redis中的数据
redis 的结构:
  	countMessages
       {chatId}:{userId}:{minutes}  count
	countInvites
       {chatId}:{userId}:{minutes}  count
	countJoins:
		{chatId}:{minutes}	count
	countLeaves
		{chatId}:{minutes}	count
*/

var (
	StatPrefixs = map[int]string{
		model.StatTypeMessageCount: "countMessages",
		model.StatTypeInviteCount:  "countInvites",
		model.StatTypeJoinChat:     "countJoins",
		model.StatTypeLeaveChat:    "countLeaves",
	}
	incStatScript = redis.NewScript(`
		local val = tonumber(redis.call('HGET', KEYS[1], KEYS[2]) or 0);
		val = val + tonumber(ARGV[1]);
		redis.call('HSET', KEYS[1], KEYS[2], val);
		return val;
	`)
)

// 统计群组消息
func StatChatMessage(chatId, userId, timestamp int64) {
	_, err := incStatCount(&model.StatCount{
		ChatId:   chatId,
		StatType: model.StatTypeMessageCount,
		UserId:   userId,
		Ts:       timestamp,
		Count:    1,
	})
	if err != nil {
		logger.Error().Err(err).Msg("stat chat message failed")
	}
}

// 进群入库
func StatsNewMembers(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	if chat == nil {
		logger.Err(err).Msg("chat is nil")
		return
	}
	chatId := chat.ID
	// 1. 创建group
	saveChatGroup(chat)
	newMembers := msg.NewChatMembers
	for _, member := range newMembers {
		userId := member.ID
		//  创建用户
		saveUser(&member)
		// 创建/更新 user-chat 关系 createOrUpdate
		UpdateChatMember(chatId, userId, "member")
		// 创建 user action
		SaveUserAction(userId, chatId, model.UserJoin)
	}
}

// 离群入库
func StatsLeave(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	if chat == nil {
		logger.Err(err).Msg("chat is nil")
		return
	}
	chatId := chat.ID
	leftMember := msg.LeftChatMember
	if leftMember == nil {
		return
	}
	userId := leftMember.ID

	RemoveChatMember(userId, chatId)
	// 创建 user action
	SaveUserAction(userId, chatId, model.UserLeft)
}

// IncStatCount 增加 redis 统计值
func incStatCount(data *model.StatCount) (int64, error) {
	prefix := StatPrefixs[data.StatType]
	minutes := data.Ts / 60

	keys := []string{prefix}
	if data.StatType == model.StatTypeInviteCount ||
		data.StatType == model.StatTypeMessageCount {
		keys = append(keys, fmt.Sprintf("%d:%d:%d", data.ChatId, data.UserId, minutes))
	} else {
		keys = append(keys, fmt.Sprintf("%d:%d", data.ChatId, minutes))
	}

	args := []interface{}{data.Count}
	return incStatScript.Run(context.Background(), rdb, keys, args).Int64()
}

// InsertMessageCountBatch 将redis中的数据批量写入数据库
func InsertMessageCountBatch(items []model.StatCount) error {
	// todo 是否需要最多每次写入200个?
	err := db.Save(items).Error
	// todo 如果写入失败, 逐个写入
	return err
}

// 查询指定时间范围内按照用户id group的结果
func FindChatCountGroupByUser(
	statType int,
	chatId int64,
	startTs, endTs int64,
	offset, limit int64) (stats []model.StatCount, err error) {
	// select user_id, sum(count) as total from stat_count sc where chat_id=xx and stat_type=xx and ts > xx and ts < xx group by user_id ;
	err = db.Raw("select user_id, sum(count) as count from stat_count sc where chat_id=? and stat_type=? and ts>? and ts<? group by user_id limit ? offset ?",
		chatId, statType, startTs, endTs, limit, offset).Scan(&stats).Error
	return
}

func FindChatCount(
	statType int,
	chatId int64,
	startTs, endTs int64,
	offset, limit int64) (stats []model.StatCount, err error) {
	db.Where("stat_type = ?", statType).Where("chat_id = ?", chatId).Find(&stats)
	return nil, nil
}
