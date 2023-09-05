package utils

import (
	"bufio"
	"io"
	"log"
	"os"
	"telegramBot/model"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

var (
	Config                  ConfigData
	SettingMenuMarkup       tgbotapi.InlineKeyboardMarkup
	StaticsMarkup           tgbotapi.InlineKeyboardMarkup
	GroupWelcomeMarkup      tgbotapi.InlineKeyboardMarkup
	InviteMenuMarkup        tgbotapi.InlineKeyboardMarkup
	ReplEnableyMenuMarkup   tgbotapi.InlineKeyboardMarkup
	ReplDisableMenuMarkup   tgbotapi.InlineKeyboardMarkup
	ProhibiteMenuMarkup     tgbotapi.InlineKeyboardMarkup
	PunishMenuMarkup        tgbotapi.InlineKeyboardMarkup
	PunishMenuMarkup2       tgbotapi.InlineKeyboardMarkup
	PunishMenuMarkup3       tgbotapi.InlineKeyboardMarkup
	PunishTimeMarkup        tgbotapi.InlineKeyboardMarkup
	MemberCheckMarkup       tgbotapi.InlineKeyboardMarkup
	UserCheckMenuMarkup     tgbotapi.InlineKeyboardMarkup
	FloodSettingMenuMarkup  tgbotapi.InlineKeyboardMarkup
	SpamSettingMenuMarkup   tgbotapi.InlineKeyboardMarkup
	DarkModelMenuMarkup     tgbotapi.InlineKeyboardMarkup
	VerifySettingMenuMarkup tgbotapi.InlineKeyboardMarkup
	ScheduleSettingMarkup   tgbotapi.InlineKeyboardMarkup
	ScheduleMsgMenuMarkup   tgbotapi.InlineKeyboardMarkup
	PermissionMenuMarkup    tgbotapi.InlineKeyboardMarkup
	DeleteNotifyMenuMarkup  tgbotapi.InlineKeyboardMarkup

	ActionMap = map[model.PunishType]string{
		model.PunishTypeWarning:    "警告",
		model.PunishTypeBan:        "禁言",
		model.PunishTypeKick:       "踢出",
		model.PunishTypeBanAndKick: "踢出+封禁",
		model.PunishTypeRevoke:     "仅撤回消息+不惩罚",
	}
	NotifyTimeMap = map[model.BanTimeType]string{
		model.BanTimeType1: "10秒",
		model.BanTimeType2: "60秒",
		model.BanTimeType3: "5分钟",
		model.BanTimeType4: "30分钟",
		model.BanTimeType5: "不删除",
		model.BanTimeType6: "不提醒",
	}
	GroupInfo model.GroupInfo = model.GroupInfo{
		GroupId:   -1001480000000,
		GroupName: "流量工程群组",
		GroupType: "超级群组",
		Uid:       6450102772,
	}
)

func InitConfig() {
	// 读取配置文件
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 解析 YAML 配置文件
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
	log.Println("配置文件加载成功...:")

	_, err = tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		panic(err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv(BOT_DEBUG) == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Println("日志模块初始化成功...")
}

func (r RequestFile) NeedsUpload() bool {
	return true
}

func (r RequestFile) UploadData() (name string, ioOut io.Reader, err error) {
	file, err := os.Open(r.FileName)
	return r.FileName, bufio.NewReader(file), err
}

func (r RequestFile) SendData() string {
	return "ok"
}

func CurrentTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006/01/02 15:04")
}
