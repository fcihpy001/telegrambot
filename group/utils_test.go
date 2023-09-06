package group

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestConvertToSeconds(t *testing.T) {
	var ss []string = []string{
		"4",
		" 4 ",
		"4m",
		"4 m",
		"4min",
		" 4  min",
		"4minutes",

		"4h",
		"4hour",
		"4hours",

		"4d",
		"4day",
		"4Days",

		"4w",
		"4Week",
		"4WEEKS",

		"4month",
		"4Months",

		"4y",
		"4Year",
		"4Years",
	}
	var results = []int64{
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,
		4 * 60,

		4 * 3600,
		4 * 3600,
		4 * 3600,

		4 * 86400,
		4 * 86400,
		4 * 86400,

		4 * 86400 * 7,
		4 * 86400 * 7,
		4 * 86400 * 7,

		4 * 86400 * 30,
		4 * 86400 * 30,

		4 * 86400 * 365,
		4 * 86400 * 365,
		4 * 86400 * 365,
	}
	for i, s := range ss {
		seconds := convertToSeconds(s)
		assert.Equal(t, results[i], seconds)
	}
}

func TestParseTime(t *testing.T) {
	ss := []struct {
		str string
		tm  int64
	}{
		{"2023-08-25 14:30:00", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{" 2023-08-25   14:30:00 ", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{"2023-8-25 14:30:00", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{"2023-8-25 14:30", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		// {"14:30:00", time.Date(2023, 8, 25, 14, 30, 0, 0, time.Local).Unix()},
		{"2023-08-25", time.Date(2023, 8, 25, 0, 0, 0, 0, time.Local).Unix()},
	}

	for _, item := range ss {
		val, err := parseTime(item.str, true)
		assert.Nil(t, err)
		assert.Equal(t, item.tm, val)
	}
}

func TestSendFile(t *testing.T) {
	chatId := -1001916451498
	fn := "202309022219117234.csv"
	pth := "../doc/202309022219117234.csv"

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	mgr := GroupManager{bot}
	err = mgr.sendFile(int64(chatId), fn, pth)
	assert.Nil(t, err)
}

func TestSendSwitch(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	chatId := int64(-1001916451498)
	msg := tgbotapi.NewMessage(chatId, "test inline")
	row := tgbotapi.NewInlineKeyboardRow()
	// btn := tgbotapi.
	btn := tgbotapi.NewInlineKeyboardButtonURL("管理", fmt.Sprintf("https://t.me/goat2023_bot?start=%d", chatId))
	row = append(row, btn)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	assert.Nil(t, err)
}

func TestWebapp(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	chatId := int64(-1001916451498)
	// msg := MessageEx{
	// 	Message: tgbotapi.Message{
	// 		Chat: &tgbotapi.Chat{ID: chatId},
	// 		Text: "webapp test",
	// 	},
	// } // (chatId, "test inline")
	// row := tgbotapi.NewInlineKeyboardRow()
	// // btn := tgbotapi.
	// btn := InlineKeyboardButtonEx{
	// 	InlineKeyboardButton: tgbotapi.InlineKeyboardButton{Text: "webapp"},
	// 	WebApp:               &WebAppInfo{URL: "https://python-telegram-bot.org/static/webappbot"},
	// }
	// row := InlineKeyboardMarkupEx{
	// 	InlineKeyboard: [][]InlineKeyboardButtonEx{
	// 		{btn},
	// 	},
	// }
	// // row = append(row, btn)
	// msg.ReplyMarkup = &row

	// _, err = bot.Send(msg)
	buf, _ := json.Marshal(map[string]interface{}{
		"inline_keyboard": [][]interface{}{
			{
				map[string]interface{}{
					"text": "webapp",
					"web_app": map[string]string{
						"url": "https://python-telegram-bot.org/static/webappbot",
					},
				},
			},
		},
	})
	params := map[string]string{
		"chat_id":      fmt.Sprint(chatId),
		"text":         "webapp test",
		"reply_markup": string(buf), //
	}
	resp, err := bot.MakeRequest("sendMessage", params)
	assert.Nil(t, err)
	println(prettyJSON(resp))
}
