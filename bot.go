package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Notify(c *gin.Context) {
	//接收到bot 发来的消息

	//需要拿到用户的信息

	SuccessResp(c, nil)
}

func SendMsg(c *gin.Context) {

	////开始通过电报发送消息
	//if !BotSendMsg(bot) {
	//	ErrorResp(c, 600, "发送失败")
	//	return
	//}

	//模拟bot向 server notify接口发消息
	testRequest()

	//6.向用户返回成功的响应
	SuccessResp(c, nil)
}

var bot *tgbotapi.BotAPI
var err error

func SetupBot() {
	// 初始化一个bot
	bot, err = tgbotapi.NewBotAPI("6670867019:AAFZwsnxb0sAP4XMRvUmJI5Lm8l5UCEAoZQ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("bot info:机器人初始化成功!")

	// 设置 Webhook
	webhookURL := "https://your-domain.com/your-webhook-endpoint"
	resp, err := bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Webhook has been set to %s", webhookURL)
	fmt.Println("webhoot set response:", resp)

	// 设置日志级别
	bot.Debug = false

	// 创建一个更新的通道，以接收用户的回复
	updates := bot.ListenForWebhook("/" + bot.Token)
	fmt.Println("开始监听用户回复！")
	// 处理用户回复
	for update := range updates {
		if update.CallbackQuery != nil {
			fmt.Println("用户回复了！")
			// 用户点击了按钮
			callbackData := update.CallbackQuery.Data
			fmt.Println("用户回复了111！", update.CallbackQuery.Message)
			log.Printf("User clicked the button with data: %s", callbackData)

		}
		fmt.Println("用户回复了222！")
	}
}

// 发送消息函数
func sendMessage(bot *tgbotapi.BotAPI, userID int, text string) {
	message := tgbotapi.NewMessage(int64(userID), text)
	// 创建一个新的 InlineKeyboardMarkup，包含一个按钮
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
	message.ReplyMarkup = inlineKeyboard

	_, err := bot.Send(message)
	if err != nil {
		log.Println(err)
	}
}

type BotMsg struct {
	Users []int64 `json:"users"`
	Msg   string  `json:"msg"`
}

func testRequest() {
	url := "http://127.0.0.1:8088/bot/notify"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func BotSendMsg(botMsg BotMsg) bool {
	//初始化一个bot
	bot, err := tgbotapi.NewBotAPI("6670867019:AAFZwsnxb0sAP4XMRvUmJI5Lm8l5UCEAoZQ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("bot info:机器人初始化成功!")

	// 你的 Webhook URL（包括 HTTPS）
	webhookURL := "https://toplinkbot.com/bot/notify"

	// 设置 Webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Webhook has been set to %s", webhookURL)

	// 设置日志级别
	bot.Debug = false
	for _, userId := range botMsg.Users {

		// 创建一个新的 InlineKeyboardMarkup，包含一个按钮
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

		// 创建一个消息
		msg := tgbotapi.NewMessage(userId, botMsg.Msg)
		msg.ReplyMarkup = inlineKeyboard
		// 发送消息
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("发送消息失败:", err)

			return false
		}
		log.Println("消息已发送到！", userId)
	}

	//// 创建一个更新的通道，以接收用户的回复
	//updates := bot.ListenForWebhook("/" + bot.Token)
	//fmt.Println("开始监听用户回复！")
	//// 处理用户回复
	//for update := range updates {
	//	if update.CallbackQuery != nil {
	//		fmt.Println("用户回复了！")
	//		// 用户点击了按钮
	//		callbackData := update.CallbackQuery.Data
	//		fmt.Println("用户回复了111！", update.CallbackQuery.Message)
	//		log.Printf("User clicked the button with data: %s", callbackData)
	//
	//		// 可以在这里根据用户的回复执行不同的操作
	//	}
	//	fmt.Println("用户回复了222！")
	//}

	// 设置机器人的 Debug 模式，以便查看日志
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// 创建一个更新的通道，以接收用户的消息
	//updates := bot.ListenForWebhook("/" + bot.Token)

	//// 处理用户消息
	//for update := range updates {
	//	if update.Message == nil {
	//		continue
	//	}
	//
	//	// 获取用户消息
	//	userMessage := update.Message.Text
	//	userID := update.Message.Chat.ID
	//
	//	// 在这里可以根据用户的消息内容执行不同的操作
	//	// 例如，回复用户的消息
	//	reply := "您发送了消息: " + userMessage
	//
	//	// 创建一个简单的消息配置
	//	message := tgbotapi.NewMessage(userID, reply)
	//
	//	// 发送回复消息
	//	_, err := bot.Send(message)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}

	//// 创建一个更新的通道，以接收用户的回调查询
	//updates := bot.ListenForWebhook("/" + bot.Token)
	//
	//// 处理用户回调查询
	//for update := range updates {
	//	if update.CallbackQuery == nil {
	//		continue
	//	}
	//
	//	// 获取回调查询的数据
	//	callbackData := update.CallbackQuery.Data
	//
	//	// 获取用户信息
	//	userID := update.CallbackQuery.From.ID
	//
	//	// 在这里可以根据用户的回调查询数据执行不同的操作
	//	// 例如，根据不同的按钮执行不同的操作
	//	switch callbackData {
	//	case "button_click_1":
	//		// 用户点击了按钮1
	//		reply := "您点击了按钮1"
	//		sendMessage(bot, userID, reply)
	//
	//	case "button_click_2":
	//		// 用户点击了按钮2
	//		reply := "您点击了按钮2"
	//		sendMessage(bot, userID, reply)
	//
	//	// 添加更多的按钮处理逻辑...
	//
	//	default:
	//		// 默认操作
	//		reply := "未知按钮点击事件"
	//		sendMessage(bot, userID, reply)
	//	}
	//}

	// 设置轮询间隔（例如，每隔5秒检查一次新消息）
	pollingInterval := 5 * time.Second

	// 配置轮询参数
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0

	for {
		// 获取新消息
		updates, err := bot.GetUpdates(u)
		if err != nil {
			log.Println("Error fetching updates:", err)
			time.Sleep(pollingInterval)
			continue
		}

		// 处理每条消息
		for _, update := range updates {
			if update.Message != nil {
				// 处理用户的消息
				userMessage := update.Message.Text
				chatID := update.Message.Chat.ID

				// 在这里执行逻辑，例如回复用户的消息
				replyMessage := fmt.Sprintf("您发送的消息是：%s", userMessage)
				msg := tgbotapi.NewMessage(chatID, replyMessage)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message:", err)
				}
			}

			// 更新更新ID，以便下次轮询时不会重复处理相同的消息
			u.Offset = update.UpdateID + 1
		}

		// 休眠一段时间后再次轮询
		time.Sleep(pollingInterval)
	}
	return true
}
