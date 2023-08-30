package group

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegramBot/model"
	"telegramBot/utils"
)

func (mgr *GroupManager) welcomeSetting(update *tgbotapi.Update) {
	btn11 := model.ButtonInfo{
		Text:    "是否启用",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    "✅启用",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn13 := model.ButtonInfo{
		Text:    "关闭",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "删除上条消息",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    "✅删除",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn23 := model.ButtonInfo{
		Text:    "不删",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "🦁自定义欢迎内容🦁",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "🦚文本内容",
		Data:    "group_welcome_text",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "🍇媒体图片",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn43 := model.ButtonInfo{
		Text:    "🍵链接按钮",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "🏠返回",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn11, btn12, btn13}
	row2 := []model.ButtonInfo{btn21, btn22, btn23}
	row3 := []model.ButtonInfo{btn31}
	row4 := []model.ButtonInfo{btn41, btn42, btn43}
	row5 := []model.ButtonInfo{btn51}
	rows := [][]model.ButtonInfo{row1, row2, row3, row4, row5}
	keyboard := utils.MakeKeyboard(rows)
	msg := "🎉 进群欢迎\n\n当前状态：关闭 ❌\n删除上条消息：✅\n\n自定义欢迎内容：\n┌📸 媒体图片:❌\n├🔠 链接按钮:❌\n└📄 文本内容:❌"
	utils.SendMenu(update.CallbackQuery.Message.Chat.ID, msg, keyboard, mgr.bot)
}

func (mgr *GroupManager) welcomeTextSetting(update *tgbotapi.Update) {
	//btn11 := model.ButtonInfo{
	//	Text:    "返回",
	//	Data:    "toast",
	//	BtnType: model.BtnTypeData,
	//}
	//row := []model.ButtonInfo{btn11}
	//rows := [][]model.ButtonInfo{row}
	//keyboard := utils.MakeKeyboard(rows)
	//msg := "👉 输入内容设置你的欢迎内容，仅支持文字和emoji:"
	//utils.SendMenu(update.CallbackQuery.Message.Chat.ID, msg, keyboard, mgr.bot)

	btn := tgbotapi.NewKeyboardButton("请输入欢迎内容")
	keybord := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(btn),
	)

	message := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "不知道写什么")
	message.ReplyMarkup = keybord
	_, err := mgr.bot.Send(message)
	if err != nil {
		log.Println(err)
	}

}
