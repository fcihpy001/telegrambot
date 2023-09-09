package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var permission_selectInfo model.SelectInfo

func PermissionHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	var buttons []model.ButtonInfo
	utils.Json2Button("./config/permission.json", &buttons)

	var row []model.ButtonInfo
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btn := buttons[i]
		updatePermissionButtonStatus(&btn)
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

	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	temp := strings.Split(params, "&")
	text := temp[0]
	row, _ := strconv.Atoi(temp[1])

	content := "所有管理员(含匿名管理员)"
	if row == 1 {
		content = "拥有封禁权限的管理"
	} else if row == 2 {
		content = "拥有更改群组信息权限的管理"
	} else if row == 3 {
		content = "拥有添加新管理员权限的管理"
	} else if row == 4 {
		content = "仅群主"
	}

	//取消原有选中
	utils.PermissionMenuMarkup.InlineKeyboard[permission_selectInfo.Row][0].Text = permission_selectInfo.Text
	//选中新的
	utils.PermissionMenuMarkup.InlineKeyboard[row][0].Text = "✅" + content
	//记录选中的索引
	permission_selectInfo.Row = row
	permission_selectInfo.Column = 0
	permission_selectInfo.Text = content

	//更新数据库
	utils.GroupInfo.Permission = text
	services.SaveModel(&utils.GroupInfo, utils.GroupInfo.GroupId)

	editText := tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		utils.PermissionMenuMarkup,
	)
	bot.Send(editText)
}

func updatePermissionButtonStatus(btn *model.ButtonInfo) {
	if btn.Text == "所有管理员(含匿名管理员)" && utils.GroupInfo.Permission == "all" {
		btn.Text = "✅" + btn.Text
		permission_selectInfo.Text = btn.Text
		permission_selectInfo.Row = 0
		permission_selectInfo.Column = 0

	} else if btn.Text == "拥有封禁权限的管理" && utils.GroupInfo.Permission == "ban" {
		btn.Text = "✅" + btn.Text
		permission_selectInfo.Text = btn.Text
		permission_selectInfo.Row = 1
		permission_selectInfo.Column = 0

	} else if btn.Text == "拥有更改群组信息权限的管理" && utils.GroupInfo.Permission == "modify" {
		btn.Text = "✅" + btn.Text
		permission_selectInfo.Text = btn.Text
		permission_selectInfo.Row = 2
		permission_selectInfo.Column = 0

	} else if btn.Text == "拥有添加新管理员权限的管理" && utils.GroupInfo.Permission == "add" {
		btn.Text = "✅" + btn.Text
		permission_selectInfo.Text = btn.Text
		permission_selectInfo.Row = 3
		permission_selectInfo.Column = 0

	} else if btn.Text == "仅群主" && utils.GroupInfo.Permission == "creator" {
		btn.Text = "✅" + btn.Text
		permission_selectInfo.Text = btn.Text
		permission_selectInfo.Row = 4
		permission_selectInfo.Column = 0
	}
}
