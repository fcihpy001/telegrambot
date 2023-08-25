package services

import (
	"telegramBot/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIncStatCount(t *testing.T) {
	rdb = OpenRedis("redis://localhost:6379/0")

	count, err := IncStatCount(model.StatCount{
		ChatId:    -1001916451498,
		StatType:  model.StatTypeMessageCount,
		UserId:    1091633677,
		Timestamp: time.Now().Unix(),
		Count:     1,
	})
	assert.Nil(t, err)
	t.Log(count)
}
