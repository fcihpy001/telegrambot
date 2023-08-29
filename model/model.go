package model

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
	ID        uint64 `gorm:"id"`
	UserId    int64
	FirstName string `gorm:"type:varchar(30)"`
	Username  string `gorm:"type:varchar(30)"`
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
