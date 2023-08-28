package bot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"strconv"
	"telegramBot/group"
	"telegramBot/utils"
)

func StartBot() {

	bot, err := createBot(utils.Config.Token, os.Getenv(utils.BOT_DEBUG) == "true")
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", bot.bot.Self.UserName)

	bot.setupBotWithPool()
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

	for update := range updatesChannel {
		s, _ := json.Marshal(update)
		log.Println("update:", string(s))
		if update.Message != nil && update.Message.IsCommand() {
			bot.handleCommand(update)
		} else if update.Message != nil && update.Message.NewChatMembers != nil {
			bot.SendText(update.Message.Chat.ID, "欢迎进群。。。")
			group.GroupHandlerMessage(update.Message, bot.bot)
		} else if update.Message != nil && update.Message.LeftChatMember != nil {
			bot.SendText(update.Message.Chat.ID, "有人离开了群组。。。")
			group.GroupHandlerMessage(update.Message, bot.bot)
		} else if update.Message != nil {
			bot.handleMessage(&update)
		} else if update.CallbackQuery != nil {
			bot.handleQuery(&update)
		} else if update.InlineQuery != nil {
			fmt.Println("inline query")
		} else {
			bot.SendText(update.Message.Chat.ID, "这个问题，暂时无法处理")
		}
	}
}

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
