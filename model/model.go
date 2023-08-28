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
	ID        uint64 `gorm:"id"`
	ChatId    int64
	StatType  int
	UserId    int64
	Timestamp int64
	Count     int
}

type User struct {
	gorm.Model
	Uid       int64
	Name      string
	ChatCount int
	GroupName string
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
