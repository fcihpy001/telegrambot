package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

// 抽奖总体设置
type LuckySetting struct {
	gorm.Model
	ChatId       int64 `gorm:"uniqueIndex"`
	StartPinned  bool
	ResultPinned bool
	DeleteToken  bool // 删除口令? 什么意思
}

const (
	LuckyTypeGeneral  = "general"  // 通用抽奖
	LuckyTypeChatJoin = "chatJoin" // 指定群组报道
	LuckyTypeInvite   = "invite"   // 邀请抽奖
	LuckyTypeHot      = "hot"      // 群活跃抽奖
	LuckyTypeFun      = "fun"      // 娱乐抽奖
	LuckyTypePoints   = "points"   // 积分抽奖
	LuckyTypeAnswer   = "answer"   // 答题抽奖

	LuckySubTypeUsers = "users" // 限制抽奖人数
	LuckySubTypeTime  = "time"

	LuckySubTypeInviteRank  = "inviteRank"  // 邀请排名 邀请排名抽奖
	LuckySubTypeInviteTimes = "inviteTimes" // 邀请次数 达到邀请人数参与随机抽奖

	LuckySubTypeHotRank  = "inviteRank"  // 邀请排名 邀请排名抽奖
	LuckySubTypeHotTimes = "inviteTimes" // 邀请次数 达到邀请人数参与随机抽奖

	LuckySubTypeFunDice     = "fruits" // dice
	LuckySubTypeFunTarget   = "fruits" // target
	LuckySubTypeFunBasket   = "fruits" // basket
	LuckySubTypeFunFootball = "fruits" // football
	LuckySubTypeFunBowl     = "fruits" // bowl
	LuckySubTypeFunFruits   = "fruits" // 水果机

	LuckyStatusStart  = "start"
	LuckyStatusEnd    = "end"
	LuckyStatusCancel = "cancel"
)

// 抽奖活动
type LuckyActivity struct {
	gorm.Model
	ChatId       int64
	UserId       int64
	LuckyName    string `gorm:"type:varchar(200)"`
	Creator      string `gorm:"type:varchar(50)"` // username
	Keyword      string `gorm:"type:varchar(100)"`
	LuckyType    string `gorm:"type:varchar(20)"`
	LuckySubType string `gorm:"type:varchar(20)"`
	LuckyCond    string // 配置信息 json
	TotalReward  string `gorm:"type:varchar(30)"`
	Status       string `gorm:"type:varchar(20)"`
	RewardDetail string // 奖励信息 json
	Results      string // 开奖信息
	StartTime    int64  // 开始时间
	EndTime      int64  // 开奖时间
	Participant  int    // 参与人数
	PartReward   int    // 中奖人数
	RewardRatio  string // 中奖率
	PushChannel  bool   // 是否推送到频道
}

func (la *LuckyActivity) ReachParticipantUsers() bool {
	if la.LuckyType == LuckyTypeGeneral && la.LuckySubType == LuckySubTypeUsers {
		if la.Participant >= la.GetLuckGeneralUsers() {
			return true
		}
	}
	return false
}

func (la *LuckyActivity) GetLuckyType() string {
	if la.LuckyType == LuckyTypeGeneral {
		if la.LuckySubType == LuckySubTypeUsers {
			return "满人开奖"
		} else if la.LuckySubType == LuckySubTypeTime {
			return "定时抽奖"
		}
	}
	return la.LuckyType + "-" + la.LuckySubType
}

func (la *LuckyActivity) GetRewards() (rewards []LuckyReward) {
	json.Unmarshal([]byte(la.RewardDetail), &rewards)
	return
}

// 满人开奖: 多少人参与后开奖
func (la *LuckyActivity) GetLuckGeneralUsers() int {
	var cond map[string]interface{}

	json.Unmarshal([]byte(la.LuckyCond), &cond)
	return int(cond["users"].(float64))
}

type LuckyRecord struct {
	gorm.Model
	LuckyId  int64
	ChatId   int64
	UserId   int64
	Username string `gorm:"type:varchar(50)"`
	Reward   string `gorm:"type:varchar(200)"` // 中奖结果
}

type LuckyReward struct {
	Name   string
	Shares int
}

type LuckyGeneral struct {
	ChatId    int64
	SubType   string // user time
	Users     int    // 限制人数
	StartTime int64
	EndTime   int64 // 到期时间
	Rewards   []LuckyReward
	Keyword   string
	Push      *bool
	Name      string // 活动名称
}
