package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuckyDataJson(t *testing.T) {
	data := LuckyData{
		Rewards: []LuckyReward{
			{"1USDT", 2},
		},
	}

	buf, err := json.Marshal(data)
	assert.Nil(t, err)
	fmt.Println("lucky data:", string(buf))

	rewards, err := json.Marshal(data.Rewards)
	assert.Nil(t, err)
	fmt.Println("rewards:", string(rewards))
}
