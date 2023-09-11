package model

import (
	"encoding/json"
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
	WarningAfterPunish  string
	BanTime             int
	DeleteNotifyMsgTime int64
	FloodSettingID      uint `gorm:"uniqueIndex"`
	SpamSettingID       uint `gorm:"uniqueIndex"`
	ProhibitedSettingID uint `gorm:"uniqueIndex"`
}

type PunishRecord struct {
	gorm.Model
	ChatId       int64
	UserId       int64
	Reason       string
	Punish       PunishType
	WarningCount int
	BanTime      int
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
	MsgCount   int64
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
	ChatId int64 `gorm:"uniqueIndex"`
	Uid    int64
	Enable bool

	StartDate     string
	EndDate       string
	StartHour     int
	EndHour       int
	Pin           bool
	RepeatHour    int
	RepeatMinute  int
	DeletePrevMsg bool
	Text          string
	Media         string
	Link          string
}

type SelectInfo struct {
	Row    int
	Column int
	Text   string
}

type GroupInfo struct {
	gorm.Model
	GroupId    int64 `gorm:";uniqueIndex"`
	GroupName  string
	GroupType  string
	Uid        int64
	Permission string
}

// 消息信息
type Message struct {
	gorm.Model
	ChatId    int64
	UserId    int64
	MessageId int64
	Timestamp int
}

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
