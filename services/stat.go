package services

import (
	"context"
	"fmt"
	"telegramBot/model"

	"github.com/redis/go-redis/v9"
)

/**
把每分钟各个群的各个用户的消息数量统计到redis中, 定时将上一分钟的统计数据写入db并清除redis中的数据
redis 的结构:
  	messageCount:{chatId}
       {userId}:{minutes}  count
	inviteCount:{chatId}
       {userId}:{minutes}  count
	joinCount:{chatId}
		{minutes}	count
	leaveCount:{chatId}
		{minutes}	count
*/

var (
	StatPrefixs = map[int]string{
		model.StatTypeMessageCount: "messageCount",
		model.StatTypeInviteCount:  "inviteCount",
		model.StatTypeJoinChat:     "joinCount",
		model.StatTypeLeaveChat:    "leaveCount",
	}
	incStatScript = redis.NewScript(`
		local val = tonumber(redis.call('HGET', KEYS[1], KEYS[2]) or 0);
		val = val + tonumber(ARGV[1]);
		redis.call('HSET', KEYS[1], KEYS[2], val);
		return val;
	`)
)

// IncStatCount 增加 redis 统计值
func IncStatCount(data model.StatCount) (int64, error) {
	prefix := StatPrefixs[data.StatType]
	minutes := data.Timestamp / 60
	var keys []string
	keys = append(keys, fmt.Sprintf("%s:%d", prefix, data.ChatId))
	if data.StatType == model.StatTypeInviteCount ||
		data.StatType == model.StatTypeMessageCount {
		keys = append(keys, fmt.Sprintf("%d:%d", data.UserId, minutes))
	} else {
		keys = append(keys, fmt.Sprintf("%d", minutes))
	}

	args := []interface{}{data.Count}
	return incStatScript.Run(context.Background(), rdb, keys, args).Int64()
}

// InsertMessageCountBatch 将redis中的数据批量写入数据库
func InsertMessageCountBatch(items []model.StatCount) {

}

// 查询指定时间范围内群聊数量
func FindChatMessageCount(
	statType int,
	chatId int64,
	startTs, endTs int64,
	offset, limit int64) ([]model.StatCount, error) {
	return nil, nil
}
