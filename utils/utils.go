package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"telegramBot/model"
	"time"
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

func TransferSecond(second int) string {
	var str string = ""
	if second == 0 {
		str = ""
	}
	if second < 60 {
		str = strconv.Itoa(second) + "秒"
	} else if second < 3600 {
		str = strconv.Itoa(second/60) + "分钟"
	} else if second < 86400 {
		str = strconv.Itoa(second/3600) + "小时"
	}
	return str
}

// 将时间字符串转换为秒
func ParseTime(time string) int {
	var count int = 0
	if strings.Contains(time, "秒") {
		arr := strings.Split(time, "秒")
		count, _ = strconv.Atoi(arr[0])

	} else if strings.Contains(time, "分钟") {
		arr := strings.Split(time, "分钟")
		count, _ = strconv.Atoi(arr[0])
		count = count * 60
	} else if strings.Contains(time, "小时") {
		arr := strings.Split(time, "小时")
		count, _ = strconv.Atoi(arr[0])
		count = count * 3600

	} else if strings.Contains(time, "不提醒") {
		count = -1
	} else if strings.Contains(time, "不删除") {
		count = 0
	}
	return count
}

func IsInDateSpan(startDateStr string, endDateStr string) bool {

	// 将日期字符串解析为时间
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		fmt.Println("无法解析开始日期:", err)
		return false
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		fmt.Println("无法解析结束日期:", err)
		return false
	}

	// 获取当前时间
	currentTime := time.Now()

	// 判断当前时间是否在两个日期之间
	if currentTime.After(startDate) && currentTime.Before(endDate) {
		return true
	}
	return false
}
