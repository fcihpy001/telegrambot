package group

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WebAppInfo struct {
	URL string `json:"url"`
}

type MessageEx struct {
	tgbotapi.Message
	ReplyMarkup *InlineKeyboardMarkupEx `json:"reply_markup,omitempty"`
}

type InlineKeyboardMarkupEx struct {
	InlineKeyboard [][]InlineKeyboardButtonEx `json:"inline_keyboard"`
}

type InlineKeyboardButtonEx struct {
	tgbotapi.InlineKeyboardButton
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

func SendTestWebApp(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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
		"chat_id":      fmt.Sprint(update.Message.Chat.ID),
		"text":         "webapp test",
		"reply_markup": string(buf), //
	}
	resp, err := bot.MakeRequest("sendMessage", params)
	if err != nil {
		logger.Err(err).Msg("send webapp failed")
	}
	println(prettyJSON(resp))
}
