package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"telegramBot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

var (
	Config                 ConfigData
	SettingMenuMarkup      tgbotapi.InlineKeyboardMarkup
	StaticsMarkup          tgbotapi.InlineKeyboardMarkup
	GroupWelcomeMarkup     tgbotapi.InlineKeyboardMarkup
	InviteMenuMarkup       tgbotapi.InlineKeyboardMarkup
	ReplEnableyMenuMarkup  tgbotapi.InlineKeyboardMarkup
	ReplDisableMenuMarkup  tgbotapi.InlineKeyboardMarkup
	ProhibiteMenuMarkup    tgbotapi.InlineKeyboardMarkup
	PunishMenuMarkup       tgbotapi.InlineKeyboardMarkup
	PunishTimeMarkup       tgbotapi.InlineKeyboardMarkup
	MemberCheckMarkup      tgbotapi.InlineKeyboardMarkup
	UserCheckMenuMarkup    tgbotapi.InlineKeyboardMarkup
	FloodSettingMenuMarkup tgbotapi.InlineKeyboardMarkup
	SpamSettingMenuMarkup  tgbotapi.InlineKeyboardMarkup

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
	ss := strings.Split(Config.Token, ":")
	botUid, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("telegram token invalid: %s %v", Config.Token, err))
	}
	Config.botUserId = botUid

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
