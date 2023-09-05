package setting

import (
	"fmt"
	"strings"
	"telegramBot/model"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Settings(chatId int64, chatType string, content string, bot *tgbotapi.BotAPI) {

	var buttons []model.ButtonInfo
	utils.Json2Button("startMenu.json", &buttons)

	var row []model.ButtonInfo
	var rows [][]model.ButtonInfo
	for i := 1; i <= len(buttons); i++ {
		btn := buttons[i-1]
		if strings.Contains(btn.Data, "群接龙") {
			btn.Data = fmt.Sprintf("https://t.me/%s?start=%d", &bot.Self.UserName, chatId)
			btn.BtnType = model.BtnTypeUrl
			//if chatType == "private" {
			//	btn.Data = fmt.Sprintf("group_solitaire?chatId=%d")
			//	btn.BtnType = model.BtnTypeData
			//}

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
	//func Settings(chatId int64, chatType string, content string, bot *tgbotapi.BotAPI) {
	//	_ = model.ButtonInfo{
	//		Text:    "🌺抽奖活动",
	//		Data:    "lucky_activity",
	//		BtnType: model.BtnTypeData,
	//	}
	//	btn12 := model.ButtonInfo{
	//		Text:    "😊专属邀请链接生成",
	//		Data:    "group_invite_setting",
	//		BtnType: model.BtnTypeData,
	//	}
	//	// 当公共群组中时, 跳转私人聊天中
	//	btn21 := model.ButtonInfo{
	//		Text:    "👨‍🎓群接龙",
	//		Data:    fmt.Sprintf("https://t.me/%s?start=%d", utils.GetBotUserName(), chatId),
	//		BtnType: model.BtnTypeUrl,
	//	}
	//	if chatType == "private" {
	//		println()
	//		btn21 = model.ButtonInfo{
	//			Text:    "👨‍🎓群接龙",
	//			Data:    fmt.Sprintf("group_solitaire?chatId=%d"),
	//			BtnType: model.BtnTypeData,
	//		}

	if len(buttons)%2 != 0 {
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.SettingMenuMarkup = keyboard

	content = fmt.Sprintf("设置【%s】群组，选择要更改的项目", utils.GroupInfo.GroupName)
	utils.SendMenu(chatId, content, keyboard, bot)
}
