package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

var Config ConfigData

func InitConfig() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// 读取配置文件
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 解析 YAML 配置文件
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
	// 解析bot userId
	ss := strings.Split(Config.Token, ":")
	uid, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil || uid == 0 {
		panic(fmt.Sprintf("invalid token: %s", Config.Token))
	}
	Config.botUserId = uid

	log.Println("配置文件加载成功...:")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if Config.Debug || os.Getenv(BOT_DEBUG) == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Println("日志模块初始化成功...")
}

func GetBotUserId() int64 {
	return Config.botUserId
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
