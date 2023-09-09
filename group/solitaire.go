package group

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 接龙

var (
	solitaireStatus = map[string]string{
		model.SolitaireStatusActive: "收集中",
		model.SolitaireStatusEnded:  "已结束",
	}
)

func SolitaireHome(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mgr := &GroupManager{bot}
	mgr.SolitaireIndex(update)
}

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
	content := fmt.Sprintf("🐉【%s】群接龙\n使用接龙来帮你方便快捷的收集用户提交的信息。\n\n", chat.Title)

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
			Text:    "🗑️",
			Data:    fmt.Sprintf("solitaire_delete?id=%d", item.ID),
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
	nextStep := "solitaire_create_last_step?" + param
	if typVal == "limitUser" {
		nextStep = "solitaire_create_step2?typ=limitUser"
	} else if typVal == "limitTime" {
		nextStep = "solitaire_create_step2?typ=limitTime"
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
				Data:    nextStep,
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

func getQueryVal(param string, key string) string {
	kvs, err := url.ParseQuery(param)
	if err != nil {
		logger.Err(err).Stack().Str("param", param).Msg("solitaire: invalid param")
		return ""
	}
	vals := kvs[key]
	if len(vals) == 0 {
		logger.Err(err).Stack().Str("param", param).Str("key", key).Msg("solitaire: not found query key")
		return ""
	}
	return vals[0]
}

// param: typ=limitTime/limitUser
func SolitaireCreateStep2(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	typ := getQueryVal(param, "typ")
	if typ != "limitTime" && typ != "limitUser" {
		logger.Error().Msgf("invalid solitaire type: %s", typ)
		return
	}
	if typ == "limitTime" {
		SolitaireCreateStep2LimitTime(update, bot, param)
	} else {
		// limitUser
		SolitaireCreateStep2LimitUser(update, bot, param)
	}
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
				Data:    "solitaire_create_limit_time?unit=hour&howmany=5",
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
				Text:    "👉最后一步",
				Data:    "solitaire_create_last_step?" + param,
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

func SolitaireCreateStep2LimitUser(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	cb := update.CallbackQuery
	msg := cb.Message
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
	// 等待用户输入 接龙规则
	StartAdminConversation(chatId, chatId, cb.From.ID, int64(msg.MessageID),
		ConversationWaitSolitaireUsers, param, nil)
}

func SolitaireCreateLastStep(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	cb := update.CallbackQuery
	msg := cb.Message
	println("SolitaireCreateLastStep param:", param)
	println(prettyJSON(update))
	chat := msg.Chat
	chatId := chat.ID
	content := "🐉创建接龙\n\n  最后一步：输入接龙规则或介绍\n"

	reply := tgbotapi.NewEditMessageText(chatId, msg.MessageID, content)
	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Msg("create solitaire last step failed")
	}
	// 等待用户输入 接龙规则
	StartAdminConversation(chatId, chatId, cb.From.ID, int64(msg.MessageID),
		ConversationWaitSolitaireDesc, param, nil)
}

func SolitaireDelete(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	msg := update.CallbackQuery.Message
	// println(prettyJSON(update))
	chat := msg.Chat
	chatId := chat.ID
	content := "🐉群接龙\n\n是否确认删除？删除后不可恢复\n"
	kvs, err := url.ParseQuery(param)
	if err != nil {
		logger.Error().Msgf("SolitaireDelete: invalid param: %v", param)
		return
	}

	if len(kvs["id"]) == 0 {
		logger.Error().Msg("SolitaireDelete: not found id to delete")
		return
	}
	sid := kvs["id"][0]
	id, err := strconv.Atoi(sid)
	if err != nil {
		logger.Err(err).Str("param", param).Msg("SolitaireDelete: invalid id")
		return
	}
	btnGroup := utils.MakeKeyboard([][]model.ButtonInfo{
		{
			{
				Text:    "🔙返回",
				Data:    "solitaire_create_limit_time_minute",
				BtnType: model.BtnTypeData,
			},
			{
				Text:    "✅确认删除",
				Data:    fmt.Sprintf("solitaire_confirm_delete?id=%d", id),
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

func SolitaireConfirmDelete(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	msg := update.CallbackQuery.Message
	// println(prettyJSON(update))
	chat := msg.Chat
	chatId := chat.ID
	content := "🐉群接龙\n\n已删除!\n"
	kvs, err := url.ParseQuery(param)
	if err != nil {
		logger.Error().Msgf("SolitaireDelete: invalid param: %v", param)
		return
	}

	if len(kvs["id"]) == 0 {
		logger.Error().Msg("SolitaireDelete: not found id to delete")
		return
	}
	sid := kvs["id"][0]
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		logger.Err(err).Str("param", param).Msg("SolitaireDelete: invalid id")
		return
	}
	services.DeleteSolitaire(id)
	btnGroup := utils.MakeKeyboard([][]model.ButtonInfo{
		{
			{
				Text:    "🔙返回",
				Data:    "solitaire_home",
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

func (mgr *GroupManager) onSolitaireCreated(update *tgbotapi.Update, sess *botConversation) {
	msg := update.Message
	userId := msg.From.ID
	chat := msg.Chat
	chatId := chat.ID
	kvs, err := url.ParseQuery(sess.data.(string))
	if err != nil {
		logger.Error().Msgf("invalid conversation data: %v", sess.data)
		return
	}

	var (
		limitUsers = 0
		limitTime  = int64(0)
	)
	if len(kvs["typ"]) == 0 {
		logger.Warn().Msg("not found solitaire param typ")
	} else {
		typ := kvs["typ"][0]
		if typ == "limitTime" {
			unit := kvs["unit"][0]
			howmany, err := strconv.ParseInt(kvs["howmany"][0], 10, 64)
			if err != nil {
				logger.Err(err).Msg("invalid solitaire param howmany")
				return
			}
			now := time.Now().Unix()
			switch unit {
			case "minute":
				limitTime = now + 60*howmany
			case "hour":
				limitTime = now + 3600*howmany
			case "day":
				limitTime = now + 86400*howmany
			default:
				logger.Error().Msg("invalid solitaire param unit")
				return
			}
		} else if typ == "limitUser" {
			users := kvs["users"][0]
			limitUsers, _ = strconv.Atoi(users)
		}
	}
	logger.Info().Msgf("create solitaire: chatId=%d userId=%d limitUsers=%d limitTime=%v",
		chatId, userId, limitUsers, limitTime)
	// message we are expected
	item, err := services.CreateSolitaire(chatId, userId, limitUsers, limitTime, msg.Text)
	if err != nil {
		logger.Err(err).Msg("create solitaire failed")
		return
	}

	// send to admin user
	btnGroup := utils.MakeKeyboard([][]model.ButtonInfo{
		{
			{
				Text:    "🔙返回",
				Data:    "solitaire_home",
				BtnType: model.BtnTypeData,
			},
		},
	})
	reply1 := tgbotapi.NewEditMessageTextAndMarkup(
		chatId,
		int(sess.messageId),
		"✅ 设置成功，点击按钮返回。",
		btnGroup)
	if _, err := mgr.bot.Send(reply1); err != nil {
		logger.Err(err).Msg("send solitaire created message to admin failed")
	}

	// send solitaire to chat group
	reply2 := tgbotapi.NewMessage(chatId, "🐉 群接龙\n\n"+msg.Text+"\n")
	reply2.ReplyMarkup = utils.MakeKeyboard([][]model.ButtonInfo{
		{
			{
				Text: "点击参加接龙",
				Data: fmt.Sprintf("https://t.me/%s?start=%s-%d",
					mgr.bot.Self.UserName, "solitaire", item.ID),
				BtnType: model.BtnTypeUrl,
			},
		},
	})
	if _, err := mgr.bot.Send(reply2); err != nil {
		logger.Err(err).Msg("send solitaire created message to group failed")
	}
}

// 用户接龙
func PlaySolitaire(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param string) {
	msg := update.Message
	chat := msg.Chat
	chatId := chat.ID
	userId := msg.From.ID
	println("PlaySolitaire chatId: userId", chatId, userId)

	ss := strings.Split(param, "-")
	if len(ss) != 2 {
		logger.Error().Msgf("PlaySolitaire: invalid param %s", param)
		return
	}
	sid, err := strconv.Atoi(ss[1])
	if err != nil {
		logger.Err(err).Msg("PlaySolitaire: invalid solitaire id")
		return
	}

	item, err := services.GetChatSolitaireById(sid)
	if err != nil {
		logger.Err(err).Msg("GetChatSolitaire failed")
		return
	}

	prevSol, err := services.GetUserSolitaire(sid, userId)
	if err != nil {
		logger.Err(err).Msg("GetUserSolitaire failed")
		return
	}
	var reply tgbotapi.Chattable
	if prevSol != nil {
		reply = tgbotapi.NewMessage(chatId,
			fmt.Sprintf("🐉群接龙\n\n%s\n\n您的接龙内容:%s\n\n输入您修改后的接龙内容:\n",
				item.Description, prevSol.Message))
	} else {
		reply = tgbotapi.NewMessage(chatId, fmt.Sprintf("🐉群接龙\n\n%s\n\n请输入您的接龙内容:\n", item.Description))
	}
	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Msg("PlaySolitaire: send message failed")
		return
	}
	// save this session, wait user's reply
	StartAdminConversation(chatId, chatId, userId, int64(msg.MessageID),
		ConversationPlaySolitaire, map[string]interface{}{
			"solitaireId":   sid,
			"prevSolitaire": prevSol,
		}, nil)
}

func (mgr *GroupManager) onSolitaireLimitUser(update *tgbotapi.Update, sess *botConversation) {
	msg := update.Message
	chat := msg.Chat
	chatId := chat.ID
	userId := msg.From.ID
	println("onSolitaireLimitUser: ", chatId, userId, sess.chatId, sess.messageId)

	param := fmt.Sprintf("typ=limitUser&users=%v", msg.Text)
	if _, err := strconv.Atoi(msg.Text); err != nil {
		logger.Err(err).Msgf("invalid solitaire limit user: %v", msg.Text)
	}

	content := "🐉创建接龙\n\n  最后一步：输入接龙规则或介绍\n"

	reply := tgbotapi.NewEditMessageText(chatId, int(sess.messageId), content)
	if _, err := mgr.bot.Send(reply); err != nil {
		logger.Err(err).Msg("create solitaire last step failed")
	}
	// 等待用户输入 接龙规则
	StartAdminConversation(sess.chatId, chatId, userId, int64(msg.MessageID),
		ConversationWaitSolitaireDesc, param, nil)
}

// 用户接龙消息的处理
func (mgr *GroupManager) onPlaySolitaireComplete(update *tgbotapi.Update, sess *botConversation) {
	msg := update.Message
	chat := msg.Chat
	chatId := chat.ID
	userId := msg.From.ID

	data := sess.data.(map[string]interface{})
	sid := data["solitaireId"].(int)
	item, err := services.GetChatSolitaireById(sid)
	if err != nil {
		logger.Err(err).Msg("GetChatSolitaire failed")
		return
	}
	// 检查该接龙是否已结束
	if item.Status != "active" {
		reply := tgbotapi.NewMessage(chatId, "接龙已结束")
		mgr.bot.Send(reply)
		return
	}

	// 用户是否接龙过
	prevSol := data["prevSolitaire"].(*model.SolitaireMessage)
	services.NewChatSolitaireMessage(item.ChatId, int64(item.ID), userId, msg.Text)

	if prevSol == nil {
		services.UpdateSolitaireStatusByIncCollected(item.ID)
	}
	// 1. 发送接龙成功
	toUserMsg := tgbotapi.NewMessage(chatId, "✅   Solitaire success! \n\nIf you need to modify the Solitaire content, please go back to the group, click the [Participate in Solitaire] button again, and then send the Solitaire content to me to modify.")
	mgr.bot.Send(toUserMsg)

	// 2. 发送接龙消息到群
	solList, err := services.GetSolitaireMessageList(int64(item.ID))
	if err != nil {
		return
	}
	content := fmt.Sprintf("🐉群接龙\n\n%s\n\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, item.Description))
	for idx, sol := range solList {
		username := mgr.getUserName(item.ChatId, sol.UserId)
		content += fmt.Sprintf("%d\\. %s\n%s\n", idx+1,
			mentionUser(username, sol.UserId),
			tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, sol.Message))
	}

	toGroupMsg := tgbotapi.NewMessage(item.ChatId, content)
	toGroupMsg.ParseMode = tgbotapi.ModeMarkdownV2
	mgr.bot.Send(toGroupMsg)
}
