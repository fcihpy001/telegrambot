package main

import (
	"telegramBot/bot"
	"telegramBot/utils"
)

func main() {
	//=======================================================
	// 1. 读取配置文件
	utils.InitConfig()
	//=======================================================
	// 2. 启动 Bot
	bot.StartBot()
}
