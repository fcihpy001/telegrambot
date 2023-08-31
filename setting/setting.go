package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramBot/model"
	"telegramBot/utils"
)

func Settings(chatId int64, bot *tgbotapi.BotAPI) {
	btn11 := model.ButtonInfo{
		Text:    "🌺抽奖活动",
		Data:    "lucky_activity",
		BtnType: model.BtnTypeData,
	}
	btn12 := model.ButtonInfo{
		Text:    "😊专属邀请链接生成",
		Data:    "group_invite_link",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "👨‍🎓群接龙",
		Data:    "group_solitaire",
		BtnType: model.BtnTypeSwitch,
	}
	btn22 := model.ButtonInfo{
		Text:    "🧝‍群统计",
		Data:    "group_statistic",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "🐞自动回复",
		Data:    "auto_reply",
		BtnType: model.BtnTypeData,
	}
	btn32 := model.ButtonInfo{
		Text:    "🦊定时消息",
		Data:    "timing_message",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "🌳入群验证",
		Data:    "group_verification",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "进群欢迎",
		Data:    "group_welcome_setting",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "🦬反垃圾",
		Data:    "anti_spam",
		BtnType: model.BtnTypeData,
	}
	btn52 := model.ButtonInfo{
		Text:    "🌓反刷屏",
		Data:    "anti_flood",
		BtnType: model.BtnTypeData,
	}
	btn61 := model.ButtonInfo{
		Text:    "⛄️违禁词",
		Data:    "prohibited_words",
		BtnType: model.BtnTypeData,
	}
	btn62 := model.ButtonInfo{
		Text:    "🌽用户检查",
		Data:    "user_check",
		BtnType: model.BtnTypeData,
	}
	btn71 := model.ButtonInfo{
		Text:    "🌗夜晚模式",
		Data:    "night_mode",
		BtnType: model.BtnTypeData,
	}
	btn72 := model.ButtonInfo{
		Text:    "🌰新群员限制",
		Data:    "new_member_limit",
		BtnType: model.BtnTypeData,
	}

	btn91 := model.ButtonInfo{
		Text:    "🥦语言切换",
		Data:    "language_switch",
		BtnType: model.BtnTypeData,
	}
	btn92 := model.ButtonInfo{
		Text:    "🏊切换其它群",
		Data:    "switch_group",
		BtnType: model.BtnTypeData,
	}

	btnRow1 := []model.ButtonInfo{btn11, btn12}
	btnRow2 := []model.ButtonInfo{btn21, btn22}
	btnRow3 := []model.ButtonInfo{btn31, btn32}
	btnRow4 := []model.ButtonInfo{btn41, btn42}
	btnRow5 := []model.ButtonInfo{btn51, btn52}
	btnRow6 := []model.ButtonInfo{btn61, btn62}
	btnRow7 := []model.ButtonInfo{btn71, btn72}
	btnRow9 := []model.ButtonInfo{btn91, btn92}

	btns := [][]model.ButtonInfo{btnRow1, btnRow2, btnRow3, btnRow4, btnRow5, btnRow6, btnRow7, btnRow9}
	keyboard := utils.MakeKeyboard(btns)
	utils.SettingMenuMarkup = keyboard
	groupName := "流量工程"
	utils.SendMenu(chatId, fmt.Sprintf("设置【%s】群组，选择要更改的项目", groupName), keyboard, bot)
}
