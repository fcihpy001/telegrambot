package setting

import (
	"fmt"
	"strings"
	"telegramBot/model"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Settings(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	var (
		chatId   int64
		chatType string
	)
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
		chatType = update.CallbackQuery.Message.Chat.Type
	} else if update.Message != nil {
		chatId = update.Message.Chat.ID
		chatType = update.Message.Chat.Type
	}
	var buttons []model.ButtonInfo
	utils.Json2Button("./config/setting.json", &buttons)

	var row []model.ButtonInfo
	var rows [][]model.ButtonInfo
	for i := 1; i <= len(buttons); i++ {
		btn := buttons[i-1]
		if strings.Contains(btn.Data, "群接龙") {
			btn.Data = fmt.Sprintf("https://t.me/%s?start=%d", bot.Self.UserName, chatId)
			btn.BtnType = model.BtnTypeUrl
			if chatType == "private" {
				btn.Data = fmt.Sprintf("group_solitaire?chatId=%d", chatId)
				btn.BtnType = model.BtnTypeData
			}

			row = []model.ButtonInfo{btn}
			rows = append(rows, row)
			continue
		}

		row = append(row, buttons[i-1])
		if i%2 == 0 && i != 0 {
			rows = append(rows, row)
			row = []model.ButtonInfo{}
		}
	}

	if len(buttons)%2 != 0 {
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.SettingMenuMarkup = keyboard

	content := fmt.Sprintf("设置【%s】群组，选择要更改的项目", utils.GroupInfo.GroupName)
	utils.SendMenu(chatId, content, keyboard, bot)
}
