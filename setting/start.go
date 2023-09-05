package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var startInfo model.GroupInfo

func StartHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	groups, err := services.GetAllGroups()
	if err != nil {
		return
	}
	fmt.Println("startInfo", groups)

	var managerRow []model.ButtonInfo
	var rows [][]model.ButtonInfo
	for i := 1; i <= len(groups); i++ {
		btn := model.ButtonInfo{
			Text:    groups[i-1].GroupName,
			Data:    "manager_group_detail:" + groups[i-1].GroupName,
			BtnType: model.BtnTypeData,
		}
		managerRow = append(managerRow, btn)
		if i%2 == 0 && i != 0 { //每两个一组，进行换行
			rows = append(rows, managerRow)
			managerRow = []model.ButtonInfo{}
		}
	}
	if len(groups)%2 != 0 {
		rows = append(rows, managerRow)
	}

	//TODO 添加完群组后，需要将信息入库
	addBtn := model.ButtonInfo{
		Text:    "+ 添加toplink到群组 +",
		Data:    "manager_group_add",
		BtnType: model.BtnTypeData,
	}

	supportBtn1 := model.ButtonInfo{
		Text:    "抽奖推送",
		Data:    "https://t.me/+w5XtbfMx6bFlMjM1",
		BtnType: model.BtnTypeUrl,
	}
	supportBtn2 := model.ButtonInfo{
		Text:    "订阅频道",
		Data:    "https://t.me/+rkFZo-A6GFNjYTFl",
		BtnType: model.BtnTypeUrl,
	}
	supportBtn3 := model.ButtonInfo{
		Text:    "官方群组",
		Data:    "https://t.me/+vQQSVgeLNiZiNmZl",
		BtnType: model.BtnTypeUrl,
	}

	addRow := []model.ButtonInfo{addBtn}
	supportRow := []model.ButtonInfo{supportBtn1, supportBtn2, supportBtn3}
	rows = append(rows, addRow)
	rows = append(rows, supportRow)
	keyboard := utils.MakeKeyboard(rows)
	content := fmt.Sprintf("👏 欢迎使用%s，如何使用：\n                \n "+
		"•  邀请 @%s 进入群组\n •  设置为管理员\n "+
		"•  在机器人私聊中发送 /start 打开设置菜单。\n\n"+
		"/help 查看我的功能\n\n\n👉 "+
		"选择下面群组进行设置：", utils.GetBotName(), utils.GetBotUserName())
	utils.SendMenu(update.Message.Chat.ID, content, keyboard, bot)
}

func ManagerGroupHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}
	if cmd == "manager_group_add" {
		mangerGroupAdd()
	} else if cmd == "manager_group_detail" {
		managerGroupDetail(update, bot, params)
	}
}

// todo 将机器人添加到群组的逻辑
func mangerGroupAdd() {
	startInfo.GroupName = "toplink群组"
	startInfo.GroupId = -1001929237671
	startInfo.GroupType = "supergroup"

	services.SaveModel(&startInfo, startInfo.GroupId)
}

func managerGroupDetail(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}
	utils.GroupInfo.GroupName = params
	Settings(update.CallbackQuery.Message.Chat.ID, "public", update.CallbackQuery.Message.Text, bot)
}
