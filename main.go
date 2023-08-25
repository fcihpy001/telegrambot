package main

import (
	"telegramBot/bot"
	"telegramBot/services"
	"telegramBot/utils"
)

func main() {
	//=======================================================
	// 1. 读取配置文件
	utils.InitConfig()
	//=======================================================
	// 2. 初始化数据库
	services.InitDB()

	//=======================================================
	// 3. 启动 Bot
	bot.StartBot()
}
