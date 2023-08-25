package utils

import (
	"bufio"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var Config ConfigData

func InitConfig() {
	// 读取配置文件
	data, err := ioutil.ReadFile("./utils/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 解析 YAML 配置文件
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
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
