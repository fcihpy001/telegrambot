package utils

import (
	"fmt"
	"log"
	"telegramBot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendText(chatId int64, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, text)
	bot.Send(msg)
}

func SendMessage(bot *tgbotapi.BotAPI, c tgbotapi.Chattable, fmt string, args ...interface{}) {
	if _, err := bot.Send(c); err != nil {
		log.Printf(fmt, args...)
		log.Println(err)
	}
}

func SendMenu(receiver int64, msg string, keybord tgbotapi.InlineKeyboardMarkup, bot *tgbotapi.BotAPI) {
	message := tgbotapi.NewMessage(receiver, msg)
	message.ReplyMarkup = keybord
	_, err := bot.Send(message)
	if err != nil {
		log.Println(err)
	}
}

func MakeKeyboard(btns [][]model.ButtonInfo) tgbotapi.InlineKeyboardMarkup {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
	for i := 0; i < len(btns); i++ {
		row := tgbotapi.NewInlineKeyboardRow()
		for j := 0; j < len(btns[i]); j++ {
			btn := tgbotapi.NewInlineKeyboardButtonData(btns[i][j].Text, btns[i][j].Data)
			if btns[i][j].BtnType == model.BtnTypeUrl {
				btn = tgbotapi.NewInlineKeyboardButtonURL(btns[i][j].Text, btns[i][j].Data)
			} else if btns[i][j].BtnType == model.BtnTypeSwitch {
				btn = tgbotapi.NewInlineKeyboardButtonSwitch(btns[i][j].Text, btns[i][j].Data)
			}
			row = append(row, btn)
		}
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, row)
	}
	return inlineKeyboard
}

func SendReply(msgId string, bot *tgbotapi.BotAPI, isAlert bool, msg string) {
	callback := tgbotapi.NewCallback(msgId, msg)
	callback.ShowAlert = isAlert
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	}
}

func SendNotify(chatId int64, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.DisableNotification = false
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func SendEditMsgMarkup(
	chatID int64,
	messageID int,
	content string,
	replyMarkup tgbotapi.InlineKeyboardMarkup,
	bot *tgbotapi.BotAPI,
) {
	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, content, replyMarkup)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("darkModelSettingMenu", err)
	}
}
