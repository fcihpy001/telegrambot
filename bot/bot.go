package bot

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BOT_DEBUG         = "BOT_DEBUG"
	POLL_TIMEOUT      = "POLL_TIMEOUT"
	TELEGRAM_APITOKEN = "TELEGRAM_APITOKEN"
)

func CreateBot(token string, debug bool) (*SmartBot, error) {
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

func StartBot() {
	bot, err := CreateBot(os.Getenv(TELEGRAM_APITOKEN), os.Getenv(BOT_DEBUG) == "true")
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", bot.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	tmo := 600
	if os.Getenv(POLL_TIMEOUT) != "" {
		if val, err := strconv.Atoi(os.Getenv(POLL_TIMEOUT)); err == nil {
			tmo = val
		}
	}
	u.Timeout = tmo

	updates := bot.bot.GetUpdatesChan(u)

	for update := range updates {
		s, _ := json.Marshal(update)
		log.Println("update:", string(s))
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.NewChatMembers != nil {
			bot.WelcomeNewMember(update.Message)
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		bot.handleUpdate(&update)
	}
}

func (bot *SmartBot) handleUpdate(update *tgbotapi.Update) {
	msg := update.Message

	switch strings.ToLower(msg.Command()) {
	case "help":
		// msg.Text = "I understand /sayhi and /status."
	case "sayhi":
		// msg.Text = "Hi :)"
	case "status":
		// msg.Text = "I'm ok."
	case "ban":
		bot.Ban(update)
	case "unban":
		bot.UnBan(update)
	case "admin":
		bot.CheckAdmin(update)
	case "mute":
	case "unmute":
	case "stat":
		fallthrough
	case "stats":
		bot.StatsMemberMessages(update)
	case "start":
		return
	case "invite":
		bot.inviteLink(update)
	default:
		return
	}
}

func (bot *SmartBot) sendMessage(c tgbotapi.Chattable, fmt string, args ...interface{}) {
	if _, err := bot.bot.Send(c); err != nil {
		log.Printf(fmt, args...)
		log.Println(err)
	}
}
