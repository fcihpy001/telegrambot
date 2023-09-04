package group

import (
	"fmt"
	"net/url"
	"strconv"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 接龙

var (
	solitaireStatus = map[string]string{
		model.SolitaireStatusActive: "收集中",
		model.SolitaireStatusEnded:  "已结束",
	}
)

// func SolitaireHome(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
// 	mgr := &GroupManager{bot}
// 	mgr.SolitaireIndex(update)
// }

// 接龙首屏 group_solitaire
func (mgr *GroupManager) SolitaireIndex(update *tgbotapi.Update) {
	msg := update.CallbackQuery.Message
	// println(prettyJSON(update))
	chat := msg.Chat
	chatId := chat.ID
	items, err := services.GetChatSolitaireList(chatId)
	if err != nil {
		logger.Err(err).Msg("get solitaire list failed")
		return
	}
	// 	🐉【toplink】Group Solitaire
	//  Use Solitaire to help you collect information submitted by users conveniently and quickly.

	// 接龙1
	// ├收集中
	// ├创建时间：2023-09-02 21:19:44
	// ├已收集：2条
	// └规则介绍：测试接龙1
	content := fmt.Sprintf("🐉【%s】群接龙\n使用接龙来帮你方便快捷的收集用户提交的信息。\n\n", utils.GetBotUserName())

	for i, item := range items {
		content += fmt.Sprintf("接龙%d\n├%s\n├创建时间：%s\n├已收集：%d条\n└规则介绍：%s\n\n",
			i+1,
			solitaireStatus[item.Status],
			item.CreatedAt,
			item.MsgCollected,
			item.Description,
		)
	}
	rows := [][]model.ButtonInfo{}
	// buttons
	for i, item := range items {
		name := fmt.Sprintf("接龙%d", i+1)
		if item.Status == model.SolitaireStatusActive {
			name += "✅"
		} else {
			name += "❌"
		}
		btn1 := model.ButtonInfo{
			Text:    name,
			Data:    "solitaire_name",
			BtnType: model.BtnTypeData,
		}
		btn2 := model.ButtonInfo{
			Text:    "文件导出",
			Data:    "solitaire_export",
			BtnType: model.BtnTypeData,
		}
		btn3 := model.ButtonInfo{
			Text:    "消息导出",
			Data:    "solitaire_messages",
			BtnType: model.BtnTypeData,
		}
		btn4 := model.ButtonInfo{
			Text:    "删除",
			Data:    "solitaire_delete",
			BtnType: model.BtnTypeData,
		}
		rows = append(rows, []model.ButtonInfo{btn1, btn2, btn3, btn4})
	}
	rows = append(rows, []model.ButtonInfo{
		{
			Text:    "➕ 新建接龙",
			Data:    "solitaire_create_step1?typ=nolimit",
			BtnType: model.BtnTypeData,
		},
	})
	rows = append(rows, []model.ButtonInfo{
		{
			Text:    "🏠 返回首页",
			Data:    "go_setting",
			BtnType: model.BtnTypeData,
		},
	})
	keyboard := utils.MakeKeyboard(rows)
	utils.GroupWelcomeMarkup = keyboard
	reply := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = mgr.bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send solitaire index failed")
	}
}

func SolitaireCreateStep1(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	mgr := &GroupManager{bot}
	msg := update.CallbackQuery.Message
	// println("SolitaireCreateStep1: " + param)
	chat := msg.Chat
	chatId := chat.ID

	kvs, err := url.ParseQuery(param)
	if err != nil {
		logger.Err(err).Str("param", param).Msg("solitaire create: invalid param")
		return
	}

	typVal := kvs["typ"][0]
	prefixFn := func(expParam string) string {
		if expParam == typVal {
			return "✅"
		}
		return ""
	}
	btnGroup := utils.MakeKeyboard([][]model.ButtonInfo{
		{
			{
				Text:    "是否限制：",
				Data:    "solitaire_if_limit",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    prefixFn("limitUser") + "限制人数",
				Data:    "solitaire_create_step1?typ=limitUser",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    prefixFn("limitTime") + "限制时间",
				Data:    "solitaire_create_step1?typ=limitTime",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    prefixFn("nolimit") + "不限",
				Data:    "solitaire_create_step1?typ=nolimit",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    "👉下一步",
				Data:    "solitaire_create_step2?" + param,
				BtnType: model.BtnTypeData,
			},
			{
				Text:    "🔙 返回",
				Data:    "solitaire_back_create_step1",
				BtnType: model.BtnTypeData,
			},
		},
	})
	utils.GroupWelcomeMarkup = btnGroup
	reply := tgbotapi.NewEditMessageTextAndMarkup(chatId,
		msg.MessageID,
		"🐉创建接龙\n\n  1️⃣第一步：设置限制",
		btnGroup)
	mgr.bot.Send(reply)
}

func btnChoosed(expUnit string, expVal int, unit string, howmany int) string {
	if expUnit == unit && expVal == howmany {
		return fmt.Sprintf("✅%d", howmany)
	}
	return fmt.Sprint(howmany)
}

// 限制时间
func SolitaireCreateStep2LimitTime(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	msg := update.CallbackQuery.Message
	// println(prettyJSON(update))
	chat := msg.Chat
	chatId := chat.ID
	content := "🐉创建接龙\n\n  2️⃣第二步：选择多久后截止\n"
	kvs, err := url.ParseQuery(param)
	if err != nil {
		logger.Err(err).Str("param", param).Msg("solitaire create step2: invalid param")
		return
	}
	var (
		unit    string
		howmany int
	)
	if len(kvs["unit"]) > 0 {
		unit = kvs["unit"][0]
	}
	if len(kvs["howmany"]) > 0 {
		howmany, _ = strconv.Atoi(kvs["howmany"][0])
	}
	btnGroup := utils.MakeKeyboard([][]model.ButtonInfo{
		{
			{
				Text:    "【按分钟】",
				Data:    "solitaire_create_limit_time_minute",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    btnChoosed(unit, howmany, "minute", 5),
				Data:    "solitaire_create_limit_time?unit=minute&howmany=5",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "minute", 10),
				Data:    "solitaire_create_limit_time?unit=minute&howmany=10",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "minute", 20),
				Data:    "solitaire_create_limit_time?unit=minute&howmany=20",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "minute", 30),
				Data:    "solitaire_create_limit_time?unit=minute&howmany=30",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    "【按小时】",
				Data:    "solitaire_create_limit_time_hour",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    btnChoosed(unit, howmany, "hour", 1),
				Data:    "solitaire_create_limit_time?unit=hour&howmany=1",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "hour", 2),
				Data:    "solitaire_create_limit_time?unit=hour&howmany=2",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "hour", 5),
				Data:    "solitaire_create_limit_time:hour?unit=hour&howmany=5",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "hour", 12),
				Data:    "solitaire_create_limit_time?unit=hour&howmany=12",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    "【按天】",
				Data:    "solitaire_create_limit_time_day",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    btnChoosed(unit, howmany, "day", 1),
				Data:    "solitaire_create_limit_time?unit=day&howmany=1",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "day", 2),
				Data:    "solitaire_create_limit_time?unit=day&howmany=2",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "day", 3),
				Data:    "solitaire_create_limit_time?unit=day&howmany=3",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "day", 5),
				Data:    "solitaire_create_limit_time?unit=day&howmany=5",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    btnChoosed(unit, howmany, "day", 10),
				Data:    "solitaire_create_limit_time?unit=day&howmany=10",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "day", 15),
				Data:    "solitaire_create_limit_time?unit=day&howmany=15",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "day", 20),
				Data:    "solitaire_create_limit_time?unit=day&howmany=20",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    btnChoosed(unit, howmany, "day", 30),
				Data:    "solitaire_create_limit_time?unit=day&howmany=30",
				BtnType: model.BtnTypeData,
			},
		},
		{
			{
				Text:    "👉下一步",
				Data:    "solitaire_create_last_step:" + param,
				BtnType: model.BtnTypeData,
			},
		},
	})
	reply := tgbotapi.NewEditMessageTextAndMarkup(chatId,
		msg.MessageID,
		content,
		btnGroup)

	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Msg("create solitaire with limit time failed")
	}
}

func SolitaireCreateStep2LimitUser(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := update.CallbackQuery.Message
	// println(prettyJSON(update))
	chat := msg.Chat
	chatId := chat.ID
	content := "🐉创建接龙\n\n	2️⃣第二步：请输入到达多少人后截止\n"

	reply := tgbotapi.NewEditMessageText(chatId,
		msg.MessageID,
		content,
	)

	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Msg("create solitaire with limit time failed")
	}
}

func SolitaireCreateLastStep(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	msg := update.CallbackQuery.Message
	println("SolitaireCreateLastStep:", param)
	println(prettyJSON(update))
	chat := msg.Chat
	chatId := chat.ID
	content := "🐉创建接龙\n\n  最后一步：输入接龙规则或介绍\n"

	reply := tgbotapi.NewEditMessageText(chatId, msg.MessageID, content)
	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Msg("create solitaire last step failed")
	}
	// 等待用户输入 接龙规则
	adminSessions[msg.Chat.ID] = &botAdminSession{
		groupChatId: 0,
		status:      "waitInput",
	}

}
