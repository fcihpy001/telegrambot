package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatTypeMessageCount = iota + 1
	StatTypeInviteCount
	StatTypeJoinChat
	StatTypeLeaveChat

	ActionJoin = "join"
	ActionLeft = "left"

	CautionTriggerByWords = "words" // 违禁词触发警告
	CautionTriggerByName  = "name"  // 用户名检查触发警告
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
	Punishment          Punishment
}
type PunishType string

const (
	PunishTypeWarning    PunishType = "warning"
	PunishTypeBan        PunishType = "ban"
	PunishTypeKick       PunishType = "kick"
	PunishTypeBanAndKick PunishType = "banAndKick"
	PunishTypeRevoke     PunishType = "revoke"
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
	PunishType          PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	BanTime             int
	DeleteNotifyMsgTime int64
	FloodSettingID      uint `gorm:"uniqueIndex"`
	SpamSettingID       uint `gorm:"uniqueIndex"`
	ProhibitedSettingID uint `gorm:"uniqueIndex"`
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

const (
	SolitaireStatusActive = "active"
	SolitaireStatusEnded  = "ended"
)

// 接龙
type Solitaire struct {
	gorm.Model
	ChatId       int64
	LimitUsers   int
	LimitTime    int    // 截止时间
	Creator      int64  // 创建用户 id
	Description  string `gorm:""`
	MsgCollected int    // 已接龙条数
	Status       string
}

type SolitaireExported struct {
	ChatId   int64
	UserId   int64
	UserName string
	NickName string // first name + last name
	Message  string
	CreateAt time.Time
}

// 消息接龙 一个用户只能接一次 如果接了多次会覆盖上次内容
type SolitaireMessage struct {
	gorm.Model
	ChatId      int64
	SolitaireId int64
	UserId      int64
	Message     string
}

// 用户警告记录
type UserCautions struct {
	gorm.Model
	UserId       int64  `gorm:"uniqueIndex:chat_user_trigger_idx"`
	ChatId       int64  `gorm:"uniqueIndex:chat_user_trigger_idx"`
	TriggerType  string `gorm:"uniqueIndex:chat_user_trigger_idx"` // 由于何种原因触发警告
	TriggerCount int64
}

type TriggerType string

const (
	TriggerTypeWords    = "words"
	TriggerTypeUsername = "username"
)

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
	AtGroup        bool
	AtUser         bool
	EthAddr        bool
	Command        bool
	LongMsg        bool
	MsgLength      int
	LongName       bool
	NameLength     int
	Exception      string
	Punishment     Punishment
}

type FloodSetting struct {
	gorm.Model
	ChatId     int64 `gorm:"uniqueIndex"`
	Uid        int64
	Enable     bool
	MsgCount   int
	Interval   int
	DeleteMsg  bool
	Punishment Punishment
}

type DarkModelSetting struct {
	gorm.Model
	ChatId       int64 `gorm:"uniqueIndex"`
	Uid          int64
	Enable       bool
	Notify       bool
	BanType      BanType
	BanTimeStart int
	BanTimeEnd   int
}

type BanType string

const (
	BanTypeMessage BanType = "BanTypeMessage"
	BanTypeMedia   BanType = "BanTypeMedia"
)
