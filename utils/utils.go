package utils

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"telegramBot/model"
)

func Json2Button(file string, models *[]model.ButtonInfo) {
	path, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(path, &models)
	if err != nil {
		log.Panic(err)
	}
}
func Json2Button2(file string, models *[][]model.ButtonInfo) {
	path, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(path, &models)
	if err != nil {
		log.Panic(err)
	}
}

func GetTimeData() [][]string {
	var times [][]string
	time1 := []string{"10秒", "30秒", "60秒"}
	time2 := []string{"5分钟", "10分钟", "30分钟"}
	time3 := []string{"1小时", "6小时", "12小时"}
	time4 := []string{"不提醒", "不删除"}
	times = append(times, time1)
	times = append(times, time2)
	times = append(times, time3)
	times = append(times, time4)
	return times
}

func TimeStr(time int) string {
	str := ""
	if time == 0 {
		str = "不提醒"
	} else if time == -1 {
		str = "不删除"
	} else if time <= 60 {
		str = strconv.Itoa(time) + "秒"
	} else if time <= 3600 {
		str = strconv.Itoa(time/60) + "分钟"
	} else if time <= 86400 {
		str = strconv.Itoa(time/3600) + "小时"
	}
	return str
}
