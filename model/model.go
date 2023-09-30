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
	gorm.Model
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
	gorm.Model
	ChatId    int64  `gorm:"uniqueIndex"`
	Title     string `gorm:"type:varchar(30)"`
	GroupType string `gorm:"type:varchar(30)"`
	Chatname  string `gorm:"type:varchar(30)"`
	Photo     string `gorm:"type:varchar(100)"`
	Location  string `gorm:"type:varchar(50)"`
}

type ButtonInfo struct {
	Text    string  `json:"text"`
	Data    string  `json:"data"`
	BtnType BtnType `json:"btn_type"`
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
	gorm.Model
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
	ChatId        int64 `gorm:"index:uniqueIndex"`
	Uid           int64
	Enable        bool
	WelcomeType   string
	WelcomeText   string
	WelcomeMedia  string
	WelcomeButton string
	DeletePrevMsg bool
	MessageId     int
}

type InviteSetting struct {
	gorm.Model
	ChatId          int64 `gorm:"index:uniqueIndex; primaryKey"`
	Uid             int64
	Enable          bool
	AutoGenerate    bool
	Notify          bool
	ExpireDate      string
	InviteLinkLimit int
	MemberLimit     int
	InviteCount     int
}

// 邀请记录
type InviteRecord struct {
	gorm.Model
	Uid                int64 `gorm:"uniqueIndex:user_chat_idx"`
	ChatId             int64 `gorm:"uniqueIndex:user_chat_idx"`
	ChatName           string
	ChatType           string
	InviteLink         string
	InviteUserId       int64
	LinkName           string
	InviteCount        int
	ExpireDate         int
	MemberLimit        int
	CreatesJoinRequest bool
}

// 自动回复设置
type ReplySetting struct {
	gorm.Model
	ChatId          int64 `gorm:"uniqueIndex"`
	Uid             int64
	Enable          bool
	ReplyList       []Reply
	DeleteReplyTime int
}

type Reply struct {
	gorm.Model
	ChatId         int64  `gorm:"primaryKey"`
	KeyWorld       string `gorm:"type:varchar(20);primaryKey" `
	ReplyWorld     string
	MatchAll       bool
	ReplySettingID uint
}

// 违禁设置
type ProhibitedSetting struct {
	gorm.Model
	ChatId              int64 `gorm:"primaryKey;uniqueIndex"`
	Uid                 int64
	Enable              bool
	World               string `gorm:"type:longtext"`
	Punish              PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	BanTime             int
	MuteTime            int
	DeleteNotifyMsgTime int
}

type PunishType string

const (
	PunishTypeWarning    PunishType = "warning"
	PunishTypeMute       PunishType = "mute"
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
	PunishType          PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	BanTime             int
	MuteTime            int
	DeleteNotifyMsgTime int
	Reason              string
	ReasonType          int
	Content             string
}

// 惩罚记录 ReasonType 1-prohibited 2-spam 3-flood 4-usercheck
type PunishRecord struct {
	gorm.Model
	ChatId       int64 `gorm:"primaryKey"`
	UserId       int64 `gorm:"primaryKey"`
	Name         string
	ReasonType   int    `gorm:"primaryKey"`
	Reason       string `gorm:"type:varchar(20)"`
	Punish       PunishType
	WarningCount int
	MuteTime     int
	BanTime      int
}

// 新人设置
type NewMemberCheck struct {
	gorm.Model
	ChatId    int64 `gorm:"uniqueIndex"`
	Uid       int64
	Enable    bool
	DelayTime int
}

// 用户合法性检测
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
	MuteTime            int
	DeleteNotifyMsgTime int
	DeleteMsg           bool
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
	Status       string `gorm:"type:varchar(30)"`
}

type SolitaireExported struct {
	ChatId   int64
	UserId   int64
	UserName string `gorm:"type:varchar(30)"`
	NickName string `gorm:"type:varchar(30)"` // first name + last name
	Message  string `gorm:"type:varchar(2000)"`
	CreateAt time.Time
}

// 消息接龙 一个用户只能接一次 如果接了多次会覆盖上次内容
type SolitaireMessage struct {
	gorm.Model
	ChatId      int64
	SolitaireId int64
	UserId      int64
	Message     string `gorm:"type:varchar(2000)"`
}

// 用户警告记录
type UserCautions struct {
	gorm.Model
	UserId       int64  `gorm:"uniqueIndex:chat_user_trigger_idx"`
	ChatId       int64  `gorm:"uniqueIndex:chat_user_trigger_idx"`
	TriggerType  string `gorm:"uniqueIndex:chat_user_trigger_idx;type:varchar(30)"` // 由于何种原因触发警告
	TriggerCount int64
}

type TriggerType string

const (
	TriggerTypeWords    = "words"
	TriggerTypeUsername = "username"
)

// 反垃圾设置
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
	DeleteMsg      bool

	Punish              PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	MuteTime            int
	BanTime             int
	DeleteNotifyMsgTime int
}

// 反刷屏设置
type FloodSetting struct {
	gorm.Model
	ChatId              int64 `gorm:"uniqueIndex"`
	Uid                 int64
	Enable              bool
	MsgCount            int
	Interval            int
	DeleteMsg           bool
	Punish              PunishType
	WarningCount        int
	WarningAfterPunish  PunishType
	MuteTime            int
	BanTime             int
	DeleteNotifyMsgTime int
}

type DarkModelSetting struct {
	gorm.Model
	ChatId        int64 `gorm:"uniqueIndex"`
	Uid           int64
	Enable        bool
	Notify        bool
	MuteType      MutType
	MuteTimeStart int
	MuteTimeEnd   int
	OnMessageId   int
	OffMessageId  int
}

type MutType string

const (
	MuteTypeMessage MutType = "MuteMessage"
	MuteTypeMedia   MutType = "MuteMedia"
)

type VerifySetting struct {
	gorm.Model
	ChatId     int64 `gorm:"uniqueIndex"`
	Uid        int64
	Enable     bool
	VerifyType string
	VerifyTime int
	PunishType string
}

type VerifyType string

const (
	VerifyTypeButton VerifyType = "VerifyTypeButton"
	VerifyTypeMath   VerifyType = "VerifyTypeMath"
	VerifyTypeCode   VerifyType = "VerifyTypeCode"
)

type Schedule struct {
	gorm.Model
	ChatId int64 `gorm:"uniqueIndex"`
	Uid    int64
}

type ScheduleMsg struct {
	gorm.Model
	ChatId int64 `gorm:"primaryKey"`
	Uid    int64
	Enable bool

	StartDate     string
	EndDate       string
	StartHour     int
	EndHour       int
	Pin           bool
	Repeat        int
	DeletePrevMsg bool
	Text          string `gorm:":varchar(30)"`
	Media         string
	Link          string
	ExecuteTime   time.Time
	MessageId     int
}

type ScheduleDelete struct {
	gorm.Model
	ChatId     int64 `gorm:"primaryKey"`
	MessageId  int
	DeleteTime time.Time
}

type SelectInfo struct {
	Row    int
	Column int
	Text   string
}

type GroupInfo struct {
	gorm.Model
	GroupId    int64 `gorm:"primaryKey;uniqueIndex"`
	GroupName  string
	GroupType  string
	Uid        int64
	Permission string
	GroupAdmin string
}

type Member struct {
	GroupId                 int64
	UserId                  int64
	Role                    string
	IsBot                   bool
	FirstName               string
	LastName                string
	UserName                string
	LanguageCode            string
	CanJoinGroups           bool
	CanReadAllGroupMessages bool
}

// 消息信息
type Message struct {
	gorm.Model
	ChatId int64
	UserId int64
	//MessageId int
	Timestamp int
}

type Task struct {
	gorm.Model
	ChatId        int64
	MessageId     int
	Type          string
	OperationTime time.Time
}
