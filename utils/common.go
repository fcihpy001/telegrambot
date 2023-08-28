package utils

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

var Config ConfigData

func InitConfig() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// 读取配置文件
	data, err := os.ReadFile("./utils/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 解析 YAML 配置文件
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if Config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Println("配置文件加载成功...:")
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
