package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	SetupBot()

	//BotSendMsg(BotMsg{
	//	Users: []int64{6450102772},
	//	Msg:   "设置【测试】群组，选择要更改的项目",
	//	//Msg:   "Toplink 通过电报交易机器人，在 Telegram 上构建一个快速而简单的钱包，并与 DEX 紧密集成，可实现闪电般快速的数据交换和追踪，允许用户通过电报应用进行自动化的 DEX 交易",
	//})
	initRouter()
}

func initRouter() {
	//1.创建路由
	r := gin.Default()

	//2.绑定路由规则，执行的函数
	r.GET("/bot/notify", Notify)
	r.GET("/bot/sendMsg", SendMsg)

	//需要启用https

	//3.监听端口，默认在8080
	r.Run(":8088")
}
