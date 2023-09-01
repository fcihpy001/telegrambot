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

type WelcomeSetting struct {
	gorm.Model
	ChatId        int64 `gorm:"index:user_chat_idx"`
	Uid           int64
	Enable        bool
	WelcomeType   string
	WelcomeText   string
	WelcomeMedia  string
	WelcomeButton string
	DeletePrevMsg bool
}

type InviteSetting struct {
	gorm.Model
	ChatId            int64 `gorm:"index:user_chat_idx"`
	Uid               int64
	Enable            bool
	AutoGenerate      bool
	Remind            bool
	LinkExpireTime    int64
	InviteLinkLimit   int
	InvitePeopleLimit int
	InviteCount       int
}

type Invite struct {
	gorm.Model
	Uid         int64 `gorm:"index:invite_idx"`
	InviteLink  string
	InviteCount int
}

type ReplySetting struct {
	gorm.Model
	ChatId          int64 `gorm:"uniqueIndex"`
	Uid             int64
	Enable          bool
	KeywordReply    []Reply
	DeleteReplyTime int
}

type Reply struct {
	gorm.Model
	ChatId         int64
	KeyWorld       string `gorm:"uniqueIndex"`
	ReplyWorld     string
	MatchAll       bool
	ReplySettingID uint
}

type ProhibitedSetting struct {
	gorm.Model
	ChatId              int64 `gorm:"uniqueIndex"`
	Uid                 int64
	Enable              bool
	World               string `gorm:"type:longtext"`
	Punish              PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	BanTime             int
	DeleteNotifyMsgTime int64
}
type PunishType string

const (
	PunishTypeWarning    PunishType = "PunishTypeWarning"
	PunishTypeBan        PunishType = "PunishTypeBan"
	PunishTypeKick       PunishType = "PunishTypeKick"
	PunishTypeBanAndKick PunishType = "PunishTypeBanAndKick"
	PunishTypeRevoke     PunishType = "PunishTypeRevoke"
)

type BanTimeType string

const (
	BanTimeType1 BanTimeType = "BanTimeType1"
	BanTimeType2 BanTimeType = "BanTimeType2"
	BanTimeType3 BanTimeType = "BanTimeType3"
	BanTimeType4 BanTimeType = "BanTimeType4"
	BanTimeType5 BanTimeType = "BanTimeType5"
	BanTimeType6 BanTimeType = "BanTimeType6"
)

type Punishment struct {
	gorm.Model
	Punish              PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	BanTime             int
	DeleteNotifyMsgTime int64
	FloodSettingID      uint
	SpamSettingID       uint
}

type NewMemberCheck struct {
	gorm.Model
	ChatId    int64 `gorm:"uniqueIndex"`
	Uid       int64
	Enable    bool
	DelayTime int
}

type UserCheck struct {
	gorm.Model
	ChatId              int64 `gorm:"uniqueIndex"`
	Uid                 int64
	NameCheck           bool
	UserNameCheck       bool
	IconCheck           bool
	SubScribe           bool
	ChannelAddr         string
	BlackUserList       string `gorm:"type:longtext"`
	Punish              PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	BanTime             int
	DeleteNotifyMsgTime int64
}

type SpamSetting struct {
	gorm.Model
	ChatId int64 `gorm:"uniqueIndex"`
	Uid    int64

	EnableAi       bool
	DDos           bool
	BlackUser      bool
	Link           bool
	ChannelCopy    bool
	ChannelForward bool
	UserForward    bool
	AtGroupId      bool
	AtUserId       bool
	EthAddr        bool
	Command        bool
	LongMsg        bool
	MsgLength      int
	LongName       bool
	NameLength     int
	PunishInfo     Punishment
}

type FloodSetting struct {
	gorm.Model
	ChatId     int64 `gorm:"uniqueIndex"`
	Uid        int64
	Enable     bool
	MsgCount   int
	Interval   int
	DeleteMsg  bool
	PunishInfo Punishment
}
