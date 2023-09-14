package model

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	LuckyEndTypeByUsers = "users" // 限制抽奖人数
	LuckyEndTypeByTime  = "time"

	LuckySubTypeInviteRank  = "inviteRank"  // 邀请排名 邀请排名抽奖
	LuckySubTypeInviteTimes = "inviteTimes" // 邀请次数 达到邀请人数参与随机抽奖

	LuckySubTypeHotRank  = "inviteRank"  // 邀请排名 邀请排名抽奖
	LuckySubTypeHotTimes = "inviteTimes" // 邀请次数 达到邀请人数参与随机抽奖

	LuckyInviteByLink = "link" // 专属邀请链接
	LuckyInviteByPull = "pull" // 拉人邀请

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
	LuckyEndType string `gorm:"type:varchar(20)"`
	LuckyCond    string // 配置信息 json
	TotalReward  string `gorm:"type:varchar(30)"`
	Status       string `gorm:"type:varchar(20)"`
	RewardDetail string // 奖励信息 json []
	Results      string // 开奖信息
	StartTime    int64  // 开始时间
	EndTime      int64  // 开奖时间
	Participant  int    // 参与人数
	PartReward   int    // 中奖人数
	RewardRatio  string // 中奖率
	PushChannel  bool   // 是否推送到频道
}

func (la *LuckyActivity) ReachParticipantUsers() bool {
	if la.LuckyEndType == LuckyEndTypeByUsers &&
		la.Participant >= la.GetLuckGeneralUsers() {
		return true
	}
	return false
}

func (la *LuckyActivity) GetLuckyType() string {
	if la.LuckyType == LuckyTypeGeneral {
		if la.LuckyEndType == LuckyEndTypeByUsers {
			return "满人开奖"
		} else if la.LuckyEndType == LuckyEndTypeByTime {
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

func (la *LuckyActivity) RecoverLuckyData() *LuckyData {
	var ld LuckyData
	if err := json.Unmarshal([]byte(la.LuckyCond), &ld); err != nil {
		logger.Err(err).Msg("unmarshal lucky activity cond failed")
	}
	if err := json.Unmarshal([]byte(la.RewardDetail), &ld.Rewards); err != nil {
		logger.Err(err).Msg("unmarshal lucky activity rewards failed")
	}
	ld.Name = la.LuckyName
	ld.ChatId = la.ChatId
	ld.Typ = la.LuckyType
	ld.SubType = la.LuckySubType
	ld.StartTime = la.StartTime
	ld.Username = la.Creator
	ld.UserId = la.UserId
	ld.EndType = la.LuckyEndType

	return &ld
}

type LuckyInvite struct {
	ChatId      int64
	UserId      int64
	Username    string
	Invitee     int64 // 被邀请用户
	InviteeName string
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

type LuckyData struct {
	ChatId    int64         `json:"-"`
	Typ       string        `json:"-"` // 1级大类
	SubType   string        `json:"-"` // 2级子类 user time
	EndType   string        `json:"-"` // 结束条件: 满人 时间到
	Name      string        `json:"-"` // 活动名称
	Rewards   []LuckyReward `json:"-"`
	UserId    int64         `json:"-"`
	Username  string        `json:"-"`
	StartTime int64

	InviteType     string `json:"inviteType,omitempty"`     // 邀请类型
	Users          int    `json:"users,omitempty"`          // 限制人数
	MinInviteCount int    `json:"minInviteCount,omitempty"` // 最少邀请人数限制
	EndTime        int64  `json:"endTime,omitempty"`        // 到期时间
	Keyword        string `json:"keyword,omitempty"`
	Push           *bool  `json:"push,omitempty"`
}

func (ld *LuckyData) GetTypeName() string {
	switch ld.Typ {
	case LuckyTypeGeneral:
		return "通用抽奖"
	case LuckyTypeInvite:
		if ld.SubType == LuckySubTypeInviteRank {
			return "邀请人数排名抽奖"
		} else {
			return "邀请人数排名抽奖"
		}
	}
	return ld.Typ + "-" + ld.SubType
}

func (ld *LuckyData) HowToParticiate(escape bool) (content string) {
	switch ld.Typ {
	case LuckyTypeGeneral:
		word := ld.Keyword
		if escape {
			word = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, word)
		}
		content = fmt.Sprintf("【如何参与？】在群组中回复关键词『%s』参与活动。", word)
	case LuckyTypeChatJoin:
	case LuckyTypeInvite:
		if ld.SubType == LuckySubTypeInviteRank {
			if escape {
				content = "【如何参与？】通过 /link 获得专属链接，使用 /link\\_stat 查看排名，到达开奖时间后，以该名单排名开奖。"
			} else {
				content = "【如何参与？】通过 /link 获得专属链接，使用 /link_stat 查看排名，到达开奖时间后，以该名单排名开奖。"
			}
		} else {

		}
	case LuckyTypeHot: // 群活跃抽奖
	case LuckyTypeFun: // 娱乐抽奖
	case LuckyTypePoints: // 积分抽奖
	case LuckyTypeAnswer:
	}

	return
}

func (ld *LuckyData) GetInviteType() string {
	if ld.InviteType == LuckyInviteByLink {
		return "专属链接"
	} else {
		return "添加成员"
	}
}
