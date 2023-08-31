package services

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"telegramBot/model"
	"time"

	"github.com/rs/zerolog/log"
)

// 定时轮询redis的数据, 清除过期数据并写入数据库

func doStatsRoutine(now int64) {
	ctx := context.Background()
	currentMinute := now / 60

	wg := &sync.WaitGroup{}
	wg.Add(4)

	batchInsertAndDelKeysFn := func(statType int) {
		defer wg.Done()

		key := StatPrefixs[statType]
		kvs, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			logger.Err(err).Msg("get message stats failed")
			return
		}
		expiredKeys, stats := convertToStatCounts(statType, kvs, currentMinute)
		if len(stats) > 0 {
			if err = InsertMessageCountBatch(stats); err != nil {
				logger.Warn().Err(err).Msgf("batch insert stat %s failed", key)
			}
		}
		if len(expiredKeys) > 0 {
			_, err = rdb.HDel(ctx, key, expiredKeys...).Result()
			if err != nil {
				logger.Warn().Err(err).Msg("delete expired keys failed")
			}
		}
	}

	go batchInsertAndDelKeysFn(model.StatTypeMessageCount)
	go batchInsertAndDelKeysFn(model.StatTypeInviteCount)
	go batchInsertAndDelKeysFn(model.StatTypeJoinChat)
	go batchInsertAndDelKeysFn(model.StatTypeLeaveChat)

	wg.Wait()
}

func convertToStatCounts(statType int, kv map[string]string, currentMinute int64) ([]string, []model.StatCount) {
	expiredKeys := []string{}
	stats := []model.StatCount{}
	for k, v := range kv {
		items := strings.Split(k, ":")
		count, err := strconv.Atoi(v)
		if err != nil {
			log.Warn().Msgf("invalid count: key=%s val=%s", k, v)
			expiredKeys = append(expiredKeys, k)
			continue
		}
		if len(items) != 3 && len(items) != 2 {
			log.Warn().Msgf("invalid key: %s", k)
			expiredKeys = append(expiredKeys, k)
			continue
		}

		minutes, err := strconv.ParseInt(items[len(items)-1], 10, 64)
		if err != nil {
			log.Warn().Msgf("invalid stat minutes: %s", items[len(items)-1])
			expiredKeys = append(expiredKeys, k)
			continue
		}
		if minutes >= currentMinute {
			continue
		}
		expiredKeys = append(expiredKeys, k)

		chatId, err := strconv.ParseInt(items[0], 10, 64)
		if err != nil {
			log.Warn().Msgf("invalid stat chatId: %s", items[0])
			continue
		}
		userId := int64(0)
		if len(items) == 3 {
			// chatId:userId:minutes
			userId, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				log.Warn().Msgf("invalid stat userId: %s", items[1])
				continue
			}
		}
		stats = append(stats, model.StatCount{
			ChatId:   chatId,
			StatType: statType,
			UserId:   userId,
			Ts:       minutes, // timestamp
			Count:    count,
			Day:      ToDay(minutes * 60),
		})
	}
	return expiredKeys, stats
}

// 每分钟的第30秒开始执行
func StatsRoutine(ctx context.Context) {
	now := time.Now().Unix()
	calcNextInterval := func(ts int64) time.Duration {
		if ts%60 >= 30 {
			return time.Second * time.Duration(90-ts%60)
		}
		return time.Second * time.Duration(ts%60)
	}
	tmr := time.NewTimer(calcNextInterval(now))

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("stats routine exit")
			return

		case <-tmr.C:
			now = time.Now().Unix()
			doStatsRoutine(now)
			now = time.Now().Unix()
			tmr = time.NewTimer(calcNextInterval(now))
		}
	}
}
