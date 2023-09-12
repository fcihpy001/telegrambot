package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"telegramBot/group"
	"telegramBot/setting"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(ctx context.Context) {
	bot, err := createBot(utils.Config.Token, os.Getenv(utils.BOT_DEBUG) == "true")
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s--%d-%s", bot.bot.Self.UserName, bot.bot.Self.ID, bot.bot.Self.FirstName)

	group.SetBot(bot.bot)

	go bot.setupBotWithPool()
	//	bot.setupBotWithWebhook()
}

func (bot *SmartBot) setupBotWithPool() {

	updateConfig := tgbotapi.NewUpdate(0)
	timeout := 600
	if os.Getenv(utils.POLL_TIMEOUT) != "" {
		if val, err := strconv.Atoi(os.Getenv(utils.POLL_TIMEOUT)); err == nil {
			timeout = val
		}
	}
	updateConfig.Timeout = timeout
	updatesChannel := bot.bot.GetUpdatesChan(updateConfig)

	// 机器人与用户的交互逻辑
	for update := range updatesChannel {
		s, _ := json.Marshal(update)
		log.Println("update:", string(s))
		// 统计
		group.DoStat(&update, bot.bot)

		if update.Message != nil && update.Message.IsCommand() { // 以/开头的指令消息
			bot.handleCommand(update)
			//如果不是管理或创建者，不响应命令信息
			member, _ := getMemberInfo(update.Message.Chat.ID, update.Message.From.ID, bot.bot)
			if member.IsAdministrator() || member.IsCreator() || update.Message.Chat.Type == "private" {
			}

		} else if update.Message != nil && update.Message.ReplyToMessage != nil { // 要求用户回复的消息
			bot.handleReply(&update)

		} else if update.Message != nil && update.Message.NewChatMembers != nil { // 新用户加入群组
			group.GroupHandlerMessage(update.Message, bot.bot)

		} else if update.Message != nil && update.Message.LeftChatMember != nil { // 用户离开群组，只统计

		} else if update.Message != nil { // 普通消息，要重点监控自定义关键词的处理
			bot.handleMessage(&update)

		} else if update.CallbackQuery != nil { // 按钮回调
			bot.handleQuery(&update)

		} else if update.InlineQuery != nil {
			fmt.Println("inline query")

		} else if update.ChatJoinRequest != nil { // 用户申请加入群组
			setting.UpdateInviteRecord(&update, bot.bot)

		} else {
			if update.Message != nil && update.Message.Chat != nil { // 未定义消息的处理
				// if chat is nil, panic
				bot.SendText(update.Message.Chat.ID, "这个问题，暂时无法处理")
			}
		}
	}
}

//lint:ignore U1000 ignore unused lint
func (bot *SmartBot) setupBotWithWebhook() {
	// load file
	certFile := utils.RequestFile{
		FileName: utils.Config.CertFile,
	}

	// 设置 Webhook
	wh, err := tgbotapi.NewWebhookWithCert(utils.Config.URL+bot.Token, certFile)
	if err != nil {
		log.Fatal(err)
	}
	_, err = bot.bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServeTLS("0.0.0.0:443", utils.Config.CertFile, utils.Config.KeyFile, nil)

	for update := range updates {
		log.Printf("%+v\n", update)
		// 处理来自 Telegram 的更新
		if update.Message != nil {
			//handleMessage(update)
		} else if update.CallbackQuery != nil {
			//query(update)
		}
	}
}

func createBot(token string, debug bool) (*SmartBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug

	return &SmartBot{
		Token: token,
		Debug: debug,
		bot:   bot,
	}, nil
}

func (bot *SmartBot) sendMessage(c tgbotapi.Chattable, fmt string, args ...interface{}) {
	if _, err := bot.bot.Send(c); err != nil {
		log.Printf(fmt, args...)
		log.Println(err)
	}
}

func (bot *SmartBot) SendText(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
