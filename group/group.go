package group

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GroupHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	query := update.CallbackQuery.Data
	switch query {
	case "group_setting":
		fmt.Println("group_setting")
	case "group_solitaire":
		fmt.Println("group_solitaire")
	case "group_record":
		fmt.Println("group_record")
	case "group_statistic":
		fmt.Println("group_statistic")
	case "group_verification":
		fmt.Println("group_verification")
	case "group_welcome":
		mgr.welcomeNewMember(update.Message)
	}
}

func GroupHandlerMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	mgr := GroupManager{
		bot: bot,
	}
	mgr.welcomeNewMember(message)
}
