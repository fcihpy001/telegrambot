package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	BOT_DEBUG         = "BOT_DEBUG"
	POLL_TIMEOUT      = "POLL_TIMEOUT"
	TELEGRAM_APITOKEN = "TELEGRAM_APITOKEN"
)

func mybot() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv(TELEGRAM_APITOKEN))
	if err != nil {
		panic(err)
	}

	bot.Debug = os.Getenv(BOT_DEBUG) == "true"
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	tmo := 600
	if os.Getenv(POLL_TIMEOUT) != "" {
		if val, err := strconv.Atoi(os.Getenv(POLL_TIMEOUT)); err == nil {
			tmo = val
		}
	}
	u.Timeout = tmo

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		s, _ := json.Marshal(update)
		log.Println("update:", string(s))
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
