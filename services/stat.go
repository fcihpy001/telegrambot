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
		Day:      ToDay(timestamp),
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
	invitedBy := int64(0)
	if msg.From != nil {
		invitedBy = msg.From.ID
	}
	// 1. 创建group
	SaveChatGroupByChat(chat)
	newMembers := msg.NewChatMembers
	for _, member := range newMembers {
		userId := member.ID
		//  创建用户
		saveUser(&member)
		// 创建/更新 user-chat 关系 createOrUpdate
		if userId == invitedBy {
			// 无人邀请
			UpdateChatMember(chatId, userId, 0, "member", int64(msg.Date))
		} else {
			UpdateChatMember(chatId, userId, invitedBy, "member", int64(msg.Date))
		}
		// 创建 user action
		SaveUserAction(userId, chatId, model.UserJoin, int64(msg.Date))
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

	RemoveChatMember(chatId, userId)
	// 创建 user action
	SaveUserAction(userId, chatId, model.UserLeft, int64(msg.Date))
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

// 查找时间范围内message总数，用户总数
func CountTotalChatMessageUsers(
	chatId int64,
	startTs, endTs int64) (int, int) {
	// SELECT sum(count) from stat_count sc  where chat_id = -1001916451498 and  ts > 28221204 ;
	// SELECT count(DISTINCT user_id)  from stat_count sc  where chat_id = -1001916451498 and  ts > 28221204  ;
	startMin := startTs / 60
	endMin := endTs / 60
	var messages, users int

	db.Raw("SELECT sum(count) as messages from stat_count sc where chat_id=? and  ts>=? and ts<?", chatId, startMin, endMin).Scan(&messages)
	db.Raw("SELECT count(DISTINCT user_id) as users  from stat_count sc  where chat_id=? and  ts>=? and ts<?", chatId, startMin, endMin).Scan(&users)
	logger.Info().Int64("messages", int64(messages)).Int64("users", int64(users)).Msg("count total messages and users")
	return messages, users
}

// 查询指定时间范围内按照用户id group的结果
func GetMessageCountGroupByUser(
	chatId int64,
	startTs, endTs int64,
	offset, limit int64) (stats []model.StatCount, err error) {
	startMin := startTs / 60
	endMin := endTs / 60
	// select user_id, sum(count) as total from stat_count sc where chat_id=xx and stat_type=xx and ts > xx and ts < xx group by user_id ;
	err = db.Raw("select user_id, sum(count) as count from stat_count sc where chat_id=? and stat_type=? and ts>=? and ts<? group by user_id order by count DESC limit ? offset ?",
		chatId, model.StatTypeMessageCount, startMin, endMin, limit, offset).Scan(&stats).Error
	return
}

// 查找时间范围内消息数量
func GroupChatMessageByDay(chatId int64, startDay, endDay string) (stats []model.StatCount, err error) {
	err = db.Raw("select day, sum(count) as count from stat_count sc where chat_id=? and stat_type=? and day>? and day<=? group by day order by count DESC",
		chatId, model.StatTypeMessageCount, startDay, endDay).Scan(&stats).Error
	return
}

func FindChatCount(
	statType int,
	chatId int64,
	startTs, endTs int64,
	offset, limit int64) (stats []model.StatCount, err error) {
	db.Where("stat_type = ?", statType).Where("chat_id = ?", chatId).Find(&stats)
	return stats, nil
}

// 查找时间范围内用户邀请数量排行
func GroupChatInviteByUser(
	chatId int64,
	startTs, endTs int64,
	limit, offset int64) (total int64, items []model.Counter, err error) {
	err = db.Raw("select count(*) as total from user_chat uc where chat_id=?  and ts>=? and ts<=?",
		chatId, startTs, endTs).Scan(&total).Error
	if err != nil {
		return
	}

	// select  sum(1) as count, invited_by , chat_id  from user_chat sc where chat_id=-1001916451498 and ts>'20230820' and ts<='20230830' group by invited_by  , chat_id order by count desc
	err = db.Raw("select sum(1) as count, invited_by, chat_id  from user_chat sc where chat_id=? and ts>? and ts<=? group by invited_by, chat_id order by count DESC limit ? offset ?",
		chatId, startTs, endTs, limit, offset).
		Find(&items).Error
	return
}

// 查找时间范围内进群、退群数量
func CountChatJoinLeft(
	action string,
	chatId int64,
	startTs, endTs int64,
) (int, error) {
	var total int
	err := db.Raw("select count(*) as total from user_action sc where chat_id=? and action=? and ts>? and ts<=?",
		chatId, action, startTs, endTs).
		Scan(&total).Error
	return total, err
}

// 查找最近n天内进群退群数量
func GroupChatJoinLeftByDay(
	action string,
	chatId int64,
	startTs, endTs int64,
) ([]model.Counter, error) {
	var items []model.Counter
	err := db.Raw("select sum(1) as count, day from user_action sc where chat_id=? and action=? and ts>? and ts<=? group by day order by day desc",
		chatId, action, startTs, endTs).
		Find(&items).Error
	return items, err
}

// 获取一段时间内
func GetLatestJoinLeftUsers(
	chatId int64,
	startTs, endTs int64,
	limit int,
) (joinList []model.User, leftList []model.User, err error) {
	records := []struct {
		Action string
		model.User
	}{}

	err = db.Raw("SELECT ua.user_id as uid, u.username , u.first_name , u.last_name, `action`  from user_action ua join `user` u on ua.user_id =u.uid where ua.chat_id=? and ua.ts>=? and ua.ts<=? and (ua.`action`='join' or ua.`action`='left') order by ts DESC limit ?",
		chatId, startTs, endTs, limit).Find(&records).Error
	if err != nil {
		return
	}
	for _, record := range records {
		if record.Action == "join" {
			joinList = append(joinList, record.User)
		} else if record.Action == "left" {
			leftList = append(leftList, record.User)
		}
	}
	return
}
