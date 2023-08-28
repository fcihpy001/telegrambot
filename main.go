package main

import (
	"context"
	"os"
	"os/signal"
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
	ctx, cancel := context.WithCancel(context.Background())
	services.Init(ctx)

	//=======================================================
	// 3. 启动 Bot
	go bot.StartBot(ctx)

	// gracelfully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until a signal is received.
	<-c
	cancel()
}
