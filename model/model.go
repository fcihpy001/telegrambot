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
}

type User struct {
	gorm.Model
	Uid          int64
	FirstName    string `gorm:"type:varchar(30)"`
	LastName     string `gorm:"type:varchar(30)"`
	Username     string `gorm:"type:varchar(30)"`
	LanguageCode string `gorm:"type:varchar(20)"`
}

type ChatGroup struct {
	ID        uint64 `gorm:"id"`
	ChatId    int64
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
)
