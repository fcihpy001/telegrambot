package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// configBot()
	mybot()
}

func configBot() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 设置 Telegram 机器人

	bot, err := tgbotapi.NewBotAPI("6670867019:AAFZwsnxb0sAP4XMRvUmJI5Lm8l5UCEAoZQ")
	bot.Debug = false
	if err != nil {
		log.Fatal(err)
	}

	// 设置 Webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://toplinkbot.com/bot/notify"))
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/bot/notify", func(c *gin.Context) {
		var update tgbotapi.Update

		if err := c.BindJSON(&update); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, "Invalid request")
			return
		}

		// 处理来自 Telegram 的更新
		if update.Message != nil {
			// 获取用户发送的消息文本
			messageText := update.Message.Text

			// 解析发消息
			if strings.HasPrefix(messageText, "/help") {
				// 如果消息以 "/help" 开头，执行相应的处理逻辑
				reply := "这是帮助信息2..."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			} else if strings.HasPrefix(messageText, "/game") {
				reply := "来玩游戏"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			} else if strings.HasPrefix(messageText, "/settings") {
				reply := "来玩游戏"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("抽*sencond* 奖活动", "button_click"),
						tgbotapi.NewInlineKeyboardButtonSwitch("专属邀请链接生成", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("群接龙", "http://www.baidu.com"),
						tgbotapi.NewInlineKeyboardButtonData("群统计", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("自动回复", "button_click"),
						tgbotapi.NewInlineKeyboardButtonData("定时消息", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("入群验证", "button_click"),
						tgbotapi.NewInlineKeyboardButtonData("进群欢迎", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("反垃圾", "button_click"),
						tgbotapi.NewInlineKeyboardButtonData("反刷屏", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("违禁词", "button_click"),
						tgbotapi.NewInlineKeyboardButtonData("用户检查", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("夜晚模式", "button_click"),
						tgbotapi.NewInlineKeyboardButtonData("新群没限制", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("下一页", "button_click"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("语言切换", "button_click"),
					),
				)
				msg.ReplyMarkup = inlineKeyboard
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			} else {
				// 如果消息不是以 "/help" 开头，执行其他处理逻辑
				reply := "感谢您的消息，但我不明白您的请求。"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
		c.JSON(http.StatusOK, "OK")
	})

	// 启动 Gin 服务器
	if err := r.RunTLS(":"+strconv.Itoa(443), "./cert/toplinkbot_com.pem", "./cert/toplinkbot_com.key"); err != nil {
		panic(err)
	}
}
