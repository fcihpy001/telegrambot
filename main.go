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
	configBot()
	//mybot()
}

var bot *tgbotapi.BotAPI
var err error

func configBot() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 设置 Telegram 机器人
	bot, err = tgbotapi.NewBotAPI("6670867019:AAFZwsnxb0sAP4XMRvUmJI5Lm8l5UCEAoZQ")
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
			handleMessage(update)
		} else if update.CallbackQuery != nil {
			query(update)
		}
		c.JSON(http.StatusOK, "OK")
	})

	// 启动 Gin 服务器
	if err := r.RunTLS(":"+strconv.Itoa(443), "./cert/toplinkbot_com.pem", "./cert/toplinkbot_com.key"); err != nil {
		panic(err)
	}
}

func handleMessage(update tgbotapi.Update) {
	// 获取用户发送的消息文本
	messageText := update.Message.Text
	// 解析发消息
	if strings.HasPrefix(messageText, "/help") {
		// 如果消息以 "/help" 开头，执行相应的处理逻辑
		reply := "这是帮助信息..."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if strings.HasPrefix(messageText, "/stat") {
		reply := "今天活跃统计功能"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if strings.HasPrefix(messageText, "/create") {
		reply := "创建抽奖活动功能有"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if strings.HasPrefix(messageText, "/luck") {
		reply := "本群正在进行的抽奖活动功能"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else if strings.HasPrefix(messageText, "/start") {
		settings(update.Message.Chat.ID)
	} else {
		reply := "感谢您的消息，但我不明白您的请求。"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func settings(chatId int64) {
	reply := "进入设置界面功能"
	msg := tgbotapi.NewMessage(chatId, reply)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌺抽奖活动", "luckydraw"),
			tgbotapi.NewInlineKeyboardButtonData("😊专属邀请链接生成", "invite_link"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨‍🎓群接龙", "groupsolit"),
			tgbotapi.NewInlineKeyboardButtonData("🧝‍群统计", "button_click"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🐞自动回复", "btn1_click"),
			tgbotapi.NewInlineKeyboardButtonData("🦊定时消息", "btn2_click"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌳入群验证", "button_click"),
			tgbotapi.NewInlineKeyboardButtonData("进群欢迎", "button_click"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦬反垃圾", "button_click"),
			tgbotapi.NewInlineKeyboardButtonData("🌓反刷屏", "button_click"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⛄️违禁词", "button_click"),
			tgbotapi.NewInlineKeyboardButtonData("🌽用户检查", "button_click"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌗夜晚模式", "button_1"),
			tgbotapi.NewInlineKeyboardButtonData("🌰新群员限制", "groupmembersetting"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚂下一页", "button_click"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥦语言切换", "button_click"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏊切换其它群", "button_click"),
			tgbotapi.NewInlineKeyboardButtonData("🪕打开vt群", "button_click"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func query(update tgbotapi.Update) {
	// 处理行内查询
	if update.CallbackQuery == nil {
		return
	}
	// 根据回调数据判断哪个按钮被点击
	switch update.CallbackQuery.Data {
	case "luckydraw":
		lucyFun(update)
	case "createlucky":
		createlucky(update)
	case "luckyrecord":
		luckyrecord(update)
	case "luckysetting":
		luckysetting(update)
	case "invite_link":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "邀请q")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

	case "groupmembersetting":
		sendMessage(update.CallbackQuery.Message.Chat.ID, "你点击了按钮-userv！")
		//// 用户点击按钮2，给用户发送一条消息
		//msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "你点击了按钮2！")
		//_, err := bot.Send(msg)
		//if err != nil {
		//	log.Println(err)
		//}
	case "groupsolit":

		// 向用户发送一个简短的通知，确认按钮点击已处理
		answer := tgbotapi.NewCallback(update.CallbackQuery.ID, "你点击了按钮2111！")
		_, err := bot.AnswerCallbackQuery(answer)
		if err != nil {
			log.Println(err)
		}
	case "settings":
		settings(update.CallbackQuery.Message.Chat.ID)
	}
}

func sendMessage(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
