package model

import (
	"gorm.io/gorm"
)

const (
	StatTypeMessageCount = iota + 1
	StatTypeInviteCount
	StatTypeJoinChat
	StatTypeLeaveChat
)

// StatCount 统计指定范围内 StatType 数量, 精度: 分钟
type StatCount struct {
	ID       uint64 `gorm:"id"`
	ChatId   int64
	StatType int
	UserId   int64
	Ts       int64 // timestamp
	Count    int
	Day      string `gorm:"type:varchar(20)"` // 20230828
}

type User struct {
	gorm.Model
	Uid          int64 `gorm:"uniqueIndex"`
	IsBot        bool
	FirstName    string `gorm:"type:varchar(30)"`
	LastName     string `gorm:"type:varchar(30)"`
	Username     string `gorm:"type:varchar(30)"`
	LanguageCode string `gorm:"type:varchar(20)"`
}

type ChatGroup struct {
	ID        uint64 `gorm:"id"`
	ChatId    int64  `gorm:"uniqueIndex"`
	Title     string `gorm:"type:varchar(30)"`
	GroupType string `gorm:"type:varchar(30)"`
	Chatname  string `gorm:"type:varchar(30)"`
	Photo     string `gorm:"type:varchar(100)"`
	Location  string `gorm:"type:varchar(50)"`
}

type ButtonInfo struct {
	Text    string
	Data    string
	BtnType BtnType
}

type BtnType string

const (
	BtnTypeUrl    BtnType = "url"
	BtnTypeData   BtnType = "data"
	BtnTypeSwitch BtnType = "switch"

	UserJoin = "join"
	UserLeft = "left"
)

// 用户与群组关系表
type UserChat struct {
	gorm.Model
	UserId    int64  `gorm:"uniqueIndex:user_chat_idx"`
	ChatId    int64  `gorm:"uniqueIndex:user_chat_idx"`
	InvitedBy int64  // 邀请人
	Status    string `gorm:"type:varchar(20)"` // member administrator
	Ts        int64
	Day       string `gorm:"type:varchar(20)"` // 20230828
}

type Counter struct {
	Count     int
	InvitedBy int64
	UserName  string
	ChatId    int64
	Day       string
}

// 用户进群退群动作
type UserAction struct {
	gorm.Model
	Action string `gorm:"type:varchar(20)"` // join, left, subscribe, unsubscribe
	UserId int64
	ChatId int64
	Ts     int64
	Day    string `gorm:"type:varchar(20)"` // 20230828
}
