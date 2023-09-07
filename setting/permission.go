package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/utils"
)

var permission_selectInfo = model.SelectInfo{
	Row:    1,
	Column: 0,
	Text:   "所有管理员(含匿名管理员",
}

func PermissionHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	var buttons []model.ButtonInfo
	utils.Json2Button("./config/permission.json", &buttons)

	var row []model.ButtonInfo
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btn := buttons[i]
		if btn.Data == "go_setting" {
			row = []model.ButtonInfo{btn}
			rows = append(rows, row)
			continue
		}
		if i == permission_selectInfo.Row {
			btn.Text = fmt.Sprintf("✅ %s", btn.Text)
		}
		btn.Data = fmt.Sprintf("%s:%d", btn.Data, i)
		row = []model.ButtonInfo{btn}
		rows = append(rows, row)
	}
	keyboard := utils.MakeKeyboard(rows)
	utils.PermissionMenuMarkup = keyboard

	content := fmt.Sprintf("⚙️  控制权限管理\n\n你可以指定哪些管理员能够设置机器人")
	utils.SendMenu(update.CallbackQuery.Message.Chat.ID, content, keyboard, bot)
}

func PermissionSelectHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	index, _ := strconv.Atoi(params)

	fmt.Println(cmd)
	fmt.Println(index)
	utils.PermissionMenuMarkup.InlineKeyboard[permission_selectInfo.Row][permission_selectInfo.Column].Text = "❌" + permission_selectInfo.Text

	utils.PermissionMenuMarkup.InlineKeyboard[index][0].Text = "✅" + utils.PermissionMenuMarkup.InlineKeyboard[index][0].Text

	//记录选中的索引
	permission_selectInfo.Row = index

	editText := tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		utils.PermissionMenuMarkup,
	)
	bot.Send(editText)
}
