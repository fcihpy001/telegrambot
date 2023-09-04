package setting

import (
	"fmt"
	"telegramBot/model"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Settings(chatId int64, bot *tgbotapi.BotAPI) {

	var buttons []model.ButtonInfo
	utils.Json2Button("startMenu.json", &buttons)

	var row []model.ButtonInfo
	var rows [][]model.ButtonInfo
	for i := 1; i <= len(buttons); i++ {
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
