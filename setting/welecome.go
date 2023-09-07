package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var welcomeSetting model.WelcomeSetting

func WelcomeHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "welcome_setting_menu" {
		welcomeSettingMenu(update, bot)

	} else if cmd == "dark_model_status" {
		statusHandler(update, bot, params)

	}
}

func welcomeSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	_ = services.GetModelData(utils.GroupInfo.GroupId, &inviteSetting)
	inviteSetting.ChatId = utils.GroupInfo.GroupId

	var btns [][]model.ButtonInfo
	utils.Json2Button2("./config/welcome.json", &btns)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(btns); i++ {
		btnArray := btns[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArray); j++ {
			btn := btnArray[j]
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.GroupWelcomeMarkup = keyboard

	//要读取用户设置的数据
	content := updateInviteSettingMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
