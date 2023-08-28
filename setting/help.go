package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramBot/model"
	"telegramBot/utils"
)

func Help(chatId int64, bot *tgbotapi.BotAPI) {
	//TODO 获取当前群的名子
	btn11 := model.ButtonInfo{
		Text:    "🌺添加toplink到群组",
		Data:    "选择群组",
		BtnType: model.BtnTypeSwitch,
	}
	btn21 := model.ButtonInfo{
		Text:    "🌺toplink官方群组",
		Data:    "https://t.me/cesjj",
		BtnType: model.BtnTypeUrl,
	}
	btn22 := model.ButtonInfo{
		Text:    "🌺toplink频道",
		Data:    "https://t.me/+rkFZo-A6GFNjYTFl",
		BtnType: model.BtnTypeUrl,
	}
	btn31 := model.ButtonInfo{
		Text:    "🌺抽奖推送频道",
		Data:    "https://t.me/+w5XtbfMx6bFlMjM1",
		BtnType: model.BtnTypeUrl,
	}
	btn32 := model.ButtonInfo{
		Text:    "🌺toplink帮助频道",
		Data:    "https://t.me/+rkFZo-A6GFNjYTFl",
		BtnType: model.BtnTypeUrl,
	}
	row1 := []model.ButtonInfo{btn11}
	row2 := []model.ButtonInfo{btn21, btn22}
	row3 := []model.ButtonInfo{btn31, btn32}
	rows := [][]model.ButtonInfo{row1, row2, row3}
	keyboard := utils.MakeKeyboard(rows)
	utils.SendMenu(chatId, "👏 欢迎使用ToplinkBot，如何使用：\n                \n •  邀请 @toplink 进入群组\n •  设置为管理员\n •  在机器人私聊中发送 /start 打开设置菜单。\n\n/help 查看我的功能\n\n\n👉 选择下面群组进行设置：", keyboard, bot)
}
