package group

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"telegramBot/model"
	"telegramBot/services"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ConversationLuckyCreateGeneralStep1 ConversationStatus = "createGeneralStep1"
	ConversationLuckyCreateGeneralStep2 ConversationStatus = "createGeneralStep2"
	ConversationLuckyCreateGeneralStep3 ConversationStatus = "createGeneralStep3"
	ConversationLuckyCreateGeneralStep4 ConversationStatus = "createGeneralStep4" // 关键词
	ConversationLuckyCreateGeneralStep5 ConversationStatus = "createGeneralStep5" // 活动名称
)

// LuckyHandler 处理抽奖部分功能
// func LuckyHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
// 	mgr := GroupManager{
// 		bot: bot,
// 	}
// 	query := update.CallbackQuery.Data
// 	switch query {
// 	case "lucky_activity":
// 		mgr.luckyActivity(update)

// 		// case "lucky_create":
// 		// 	mgr.luckyrecord(update)
// 		// case "lucky_record":
// 		// 	mgr.luckyrecord(update)
// 	}
// }

func (mgr *GroupManager) luckyActivity(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🎁【测试】抽奖\n\n发起抽奖次数：0    \n\n已开奖：0       未开奖：0       取消：0")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📎发起抽奖活动", "lucky_create"),
			tgbotapi.NewInlineKeyboardButtonData("📪查看抽奖记录", "lucky_record"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧶设置抽奖", "lucky_setting"),
			tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
		))
	msg.ReplyMarkup = inlineKeyboard
	_, err := mgr.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func luckyRecords(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	println("luckyRecords")
	cb := update.CallbackQuery
	chat := cb.Message.Chat
	chatId := chat.ID
	sidx := param.param.Get("idx")
	idx := getIntParam(&param.param, "idx")
	if idx < 0 {
		//
		logger.Info().Msg("no prev luck record")
		return nil
	}
	/*
	   🎁创建的抽奖记录

	   bnb来抢啦
	   ├满人开奖  (2人)
	   ├参与关键词：bnb
	   ├推送至频道：❌
	   ├奖品列表：
	   ├       10bnb     ×3份

	   创建者：bigwinner
	   创建时间：2023-09-07 17:04:59
	   状态：已取消 ❌       已参与：0人

	   第1条/共6条
	*/
	recordCount := services.GetLuckyActivityCount(chatId)
	content := "🎁创建的抽奖记录\n\n"
	var keyboard tgbotapi.InlineKeyboardMarkup
	if recordCount > 0 {
		record := services.GetLuckyActivity(chatId, int(idx))

		content += buildLuckyRecord(record)
		content += fmt.Sprintf("\n第%d条/共%d条\n", idx+1, recordCount)
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️上一条", fmt.Sprintf("lucky_record?idx=idx=%d", idx-1)),
				tgbotapi.NewInlineKeyboardButtonData("下一条➡️", fmt.Sprintf("lucky_record?idx=idx=%d", idx+1)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("取消抽奖", "lucky_cancel?idx="+sidx),
				tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
			))
	} else {
		content += "没有抽奖记录\n"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("取消抽奖", "lucky_cancel?idx="+sidx),
				tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
			))
	}
	reply := tgbotapi.NewEditMessageTextAndMarkup(chatId, cb.Message.MessageID, content, keyboard)
	_, err := bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send choose lucky type failed")
	}
	return nil
}

func luckyIndex(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	println("luckyIndex")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📎发起抽奖活动", "lucky_create_index"),
			tgbotapi.NewInlineKeyboardButtonData("📪查看抽奖记录", "lucky_record"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧶设置抽奖", "luckysetting"),
			tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
		))
	// todo
	msg := tgbotapi.NewEditMessageTextAndMarkup(param.chatId, param.msgId,
		"🎁【测试】抽奖\n\n发起抽奖次数：0    \n\n已开奖：0       未开奖：0       取消：0", inlineKeyboard)

	_, err := bot.Send(msg)
	if err != nil {
		logger.Err(err).Msg("send lucky index failed")
	}
	return err
}

// 发起抽奖首页: 选择抽奖类型
func luckyCreateIndex(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	println("luckyCreateIndex")
	content := "🎁 群抽奖类型:\n\n" +
		"🔥 通用抽奖：群员在群内回复指定关键词参与抽奖\n\n" +
		"🙋‍♂️ 指定群报道抽奖：A群成员进入B群回复指定关键词参与抽奖\n\n" +
		"🪁 邀请人数抽奖：群成员用[专属链接]或[添加成员]拉人进群，到指定人数后参与抽奖\n\n" +
		"🥰 群活跃抽奖：根据活跃排名抽奖，或达到活跃度参与随机抽奖\n\n" +
		"🎰 娱乐抽奖：水果机、摇骰子、飞镖、保龄球....\n\n" +
		"选择抽奖类型：\n"
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥通用抽奖", "lucky_create?typ=general"),
			// tgbotapi.NewInlineKeyboardButtonData("📪查看抽奖记录", "lucky_record"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("🧶设置抽奖", "luckysetting"),
		// 	tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙返回", "lucky"),
		),
	)
	msg := tgbotapi.NewEditMessageTextAndMarkup(param.chatId, param.msgId,
		content, inlineKeyboard)
	_, err := bot.Send(msg)
	if err != nil {
		logger.Err(err).Msg("send choose lucky type failed")
	}
	return err
}

// 主入口
func luckyCreate(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	if len(param.param["typ"]) == 0 {
		err := errors.New("not found param typ")
		logger.Err(err).Msg("lucky create failed")
		return err
	}
	typ := param.param["typ"][0]
	switch typ {
	case "general":
		content := "🎁创建通用抽奖\n\n" +
			"通用抽奖：群员在群内回复指定关键词参与抽奖\n\n" +
			"选择开奖方式：\n"
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("满人开奖", "lucky_create_general?subType=user"),
				tgbotapi.NewInlineKeyboardButtonData("定时抽奖", "lucky_create_general?subType=time"),
			),
			// tgbotapi.NewInlineKeyboardRow(
			// 	tgbotapi.NewInlineKeyboardButtonData("🧶设置抽奖", "luckysetting"),
			// 	tgbotapi.NewInlineKeyboardButtonData("🦀返回首页", "settings"),
			// ),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙返回选择抽奖类型", "lucky_create_index"),
			),
		)
		msg := tgbotapi.NewEditMessageTextAndMarkup(param.chatId, param.msgId,
			content, inlineKeyboard)
		_, err := bot.Send(msg)
		if err != nil {
			logger.Err(err).Msg("send create lucky general failed")
		}
	default:
		logger.Error().Msgf("not implement lucky type: %v", typ)
	}
	return nil
}

// 通用抽奖
func luckyCreateGeneral(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	if len(param.param["subType"]) == 0 {
		err := errors.New("not found param subType")
		logger.Err(err).Msg("lucky create general failed")
		return err
	}
	subType := param.param["subType"][0]

	var content string
	switch subType {
	case "user":
		// 满人抽奖
		content = "🎁创建通用抽奖(/cancel 命令返回首页)\n\n" +
			"请回复参与多少人后开奖：\n\n"
	case "time":
		// 定时抽奖
		content = "🎁创建通用抽奖(/cancel 命令返回首页)\n\n" +
			"请回复参与多少人后开奖：\n\n"
	}
	reply := tgbotapi.NewEditMessageText(param.chatId, param.msgId, content)
	_, err := bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send msg failed")
	}
	data := model.LuckyGeneral{
		ChatId:  param.chatId,
		SubType: subType,
	}
	StartAdminConversation(param.chatId, param.chatId, update.CallbackQuery.From.ID, int64(param.msgId),
		ConversationLuckyCreateGeneralStep1,
		&data,
		luckyCreateGeneralSteps,
	)

	return err
}

func toggleLuckySetting(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	var (
		toggleType string
		toggleVal  bool
	)

	logger.Info().Str("param", param.data).Msg("toggle lucky setting param")

	chatSetting := services.FindChatLuckySetting(param.chatId)
	if chatSetting == nil {
		// 初始值
		chatSetting = &model.LuckySetting{
			ChatId:       param.chatId,
			StartPinned:  true,
			ResultPinned: true,
			DeleteToken:  true,
		}
	}

	if len(param.param) > 0 {
		typ := param.param["typ"]
		if len(typ) == 0 {
			err := errors.New("toggleLuckySetting: not found param toggle type")
			logger.Error().Msg("toggleLuckySetting: not found param toggle type")
			return err
		}
		toggleType = typ[0]
		val := param.param["val"]
		if len(val) == 0 {
			err := errors.New("toggleLuckySetting: not found param toggle value")
			logger.Error().Msg("toggleLuckySetting: not found param toggle value")
			return err
		}
		toggleVal = toBool(val[0])
		switch toggleType {
		case "startPin":
			chatSetting.StartPinned = toggleVal
		case "endPin":
			chatSetting.ResultPinned = toggleVal
		case "deleteToken":
			chatSetting.DeleteToken = toggleVal
		}
		// update chat lucky setting
		services.UpdateChatLuckySettings(chatSetting)
	}

	content := "⚙ 抽奖设置\n\n✅ 发布置顶：\n└ 发布抽奖消息群内置顶\n✅ 结果置顶：\n└ 中奖结果消息群内置顶\n✅ 删除口令：\n└ 5分钟后自动删除群成员参加抽奖发送的口令消息"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉发布置顶", "luckysetting?alert=xxx"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(chatSetting.StartPinned, "启用"),
				"luckysetting?typ=startPin&val=1"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(!chatSetting.StartPinned, "关闭"),
				"luckysetting?typ=startPin&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📮结果置顶", "luckysetting?alert=xxx"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(chatSetting.ResultPinned, "启用"),
				"luckysetting?typ=endPin&val=1"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(!chatSetting.ResultPinned, "关闭"),
				"luckysetting?typ=endPin&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁删除口令", "luckysetting?alert=xxx"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(chatSetting.DeleteToken, "启用"),
				"luckysetting?typ=deleteToken&val=1"),
			tgbotapi.NewInlineKeyboardButtonData(
				radioOpenEmojj(!chatSetting.DeleteToken, "关闭"),
				"luckysetting?typ=deleteToken&val=0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📡返回到抽奖", "lucky$"),
		))
	reply := tgbotapi.NewEditMessageTextAndMarkup(param.chatId, param.msgId, content, inlineKeyboard)

	_, err := bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send toggleLuckySetting message failed")
	}
	return err
}

func buildRewardContent(data *model.LuckyGeneral) string {
	return "🎁创建通用抽奖  ( /cancel 命令返回首页)\n\n" + buildRewardBody(data)
}

func buildRewardBody(data *model.LuckyGeneral) string {
	content := ""
	if data.Name != "" {
		content += data.Name + "\n"
	}
	if data.Keyword != "" {
		content += fmt.Sprintf("├参与关键词：%s\n", data.Keyword)
	}
	if data.Push != nil {
		if *data.Push {
			content += "├推送至频道：❌\n"
		} else {
			content += "├推送至频道：✅\n"
		}
	}
	content += fmt.Sprintf("├满人开奖  (%v人)\n├奖品列表:", data.Users)
	for _, reward := range data.Rewards {
		if reward.Shares > 0 {
			content += fmt.Sprintf("├       %s    x %d份\n", reward.Name, reward.Shares)
		} else {
			content += fmt.Sprintf("├       %s\n", reward.Name)
		}
	}

	return content
}

// 用于展示抽奖活动
func buildRewardInfo(data *model.LuckyGeneral) string {
	content := fmt.Sprintf("%s\n├开奖时间：%s\n├参与关键词：%s\n├奖品列表：\n",
		data.Name,
		yyyymmddhhmmss(data.StartTime),
	)
	for _, reward := range data.Rewards {
		if reward.Shares > 0 {
			content += fmt.Sprintf("├       %s    x %d份\n", reward.Name, reward.Shares)
		} else {
			content += fmt.Sprintf("├       %s\n", reward.Name)
		}
	}
	return content
}

// 满人抽奖: step1 输入人数
func luckyCreateGeneralSteps(update *tgbotapi.Update, bot *tgbotapi.BotAPI, sess *botConversation) error {
	text := update.Message.Text
	if text == "/cancel" {
		return nil
	}
	data := sess.data.(*model.LuckyGeneral)
	status := sess.status

	switch status {
	case ConversationLuckyCreateGeneralStep1:
		sess.status = ConversationLuckyCreateGeneralStep2
		users, err := strconv.Atoi(text)
		if err != nil {
			logger.Err(err).Msg("invalid user arg")
		}
		data.Users = users
		content := fmt.Sprintf("🎁创建通用抽奖  ( /cancel 命令返回首页)\n\n├满人开奖  (%s人)\n\n请回复第一个奖品的名称（如：1USDT）：", text)
		sendText(bot, update.Message.Chat.ID, content)

	case ConversationLuckyCreateGeneralStep2:
		sess.status = ConversationLuckyCreateGeneralStep3
		reward := model.LuckyReward{
			Name: text,
		}
		data.Rewards = append(data.Rewards, reward)
		content := buildRewardContent(data)
		content += "\n请输入该奖品有多少份：\n"
		sendText(bot, update.Message.Chat.ID, content)

	case ConversationLuckyCreateGeneralStep3: // 奖品多少份
		val, err := strconv.Atoi(text)
		if err != nil {
			logger.Err(err).Msgf("invalid reward shares: %v", text)
		}
		data.Rewards[len(data.Rewards)-1].Shares = val
		sess.status = ConversationLuckyCreateGeneralStep2

		// 这里可以结束进入下一步, 也可以继续添加奖品
		content := buildRewardContent(data)
		content += "\n回复奖品名称，继续添加：\n"
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("👉结束添加奖品，进入下一步👈", "lucky_create_keywords"),
			),
		)
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, content)
		reply.ReplyMarkup = inlineKeyboard
		if _, err = bot.Send(reply); err != nil {
			logger.Err(err).Msg("send lucky create keywords failed")
		}

	case ConversationLuckyCreateGeneralStep4: // 关键词
		data.Keyword = text
		content := buildRewardContent(data)
		content += "\n请输入抽奖活动名称：\n"
		sess.status = ConversationLuckyCreateGeneralStep5
		sendText(bot, update.Message.Chat.ID, content)

	case ConversationLuckyCreateGeneralStep5: // 活动名称
		data.Name = text
		content := buildRewardContent(data)
		// 是否推送
		content += "\n❓ 是否推送到 小微抽奖推送频道 ，可以推广您的群组，并让更多人参加抽奖。\n" +
			"⚠ 推送请遵守《小微抽奖推送规则》 禁止推送测试抽奖、奖品无意义的抽奖、设置要求过高无法达到条件的抽奖，违反永久禁用推送\n" +
			"===============\n" +
			"‼️‼️️抽奖推送规则调整：\n" +
			"推送的抽奖，参加抽奖成员必须先关注抽奖推送频道，不推送则没有限制\n\n" +
			"请选择是否推送到频道：\n"
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅已知晓，推送", "lucky_push?result=1"),
				tgbotapi.NewInlineKeyboardButtonData("❌不推送", "lucky_push?result=0"),
			),
		)
		reply := tgbotapi.NewMessage(sess.chatId, content)
		reply.ReplyMarkup = inlineKeyboard
		if _, err := bot.Send(reply); err != nil {
			logger.Err(err).Stack().Msg("send msg failed")
		}
	}

	return nil
}

// callback
func luckyCreateKeywords(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	cb := update.CallbackQuery
	chat := cb.Message.Chat
	chatId := chat.ID
	sess := GetConversation(chatId)
	if sess == nil {
		logger.Error().Stack().Int64("chatId", chatId).Msg("not found session")
		return errors.New("not found session")
	}
	content := buildRewardContent(sess.data.(*model.LuckyGeneral))
	content += "\n👉 请回复参与抽奖关键词：\n"
	sess.status = ConversationLuckyCreateGeneralStep4
	sendEditText(bot, chatId, cb.Message.MessageID, content)

	return nil
}

// callback 满人抽奖: 是否推送
func luckyCreatePush(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	push := param.param["result"]
	println("push:", push[0])

	cb := update.CallbackQuery
	chat := cb.Message.Chat
	chatId := chat.ID
	sess := GetConversation(chatId)
	if sess == nil {
		logger.Error().Stack().Int64("chatId", chatId).Msg("not found session")
		return errors.New("not found session")
	}
	data := sess.data.(*model.LuckyGeneral)
	pushVal := false
	data.Push = &pushVal
	content := buildRewardContent(sess.data.(*model.LuckyGeneral))
	content += "\n🥳恭喜！已完成所有内容，是否发布到群组?\n" // todo 群组名称
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅发布抽奖", "lucky_publish?result=1"),
			tgbotapi.NewInlineKeyboardButtonData("❌取消发布", "lucky_publish?result=0"),
		),
	)
	reply := tgbotapi.NewMessage(sess.chatId, content)
	reply.ReplyMarkup = inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Stack().Msg("send msg failed")
	}
	return nil
}

// callback 是否发布
func luckyCreatePublish(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
	if len(param.param["result"]) == 0 {
		logger.Error().Stack().Msg("invalid param result")
		return nil
	}
	result := toBool(param.param["result"][0])
	if !result {
		// 取消发布 返回首页
		luckyCreateIndex(update, bot, param)
		return nil
	}
	// 发布
	cb := update.CallbackQuery
	chat := cb.Message.Chat
	chatId := chat.ID
	sess := GetConversation(chatId)
	if sess == nil {
		logger.Error().Stack().Int64("chatId", chatId).Msg("not found session")
		return errors.New("not found session")
	}

	data := sess.data.(*model.LuckyGeneral)

	content := buildRewardContent(sess.data.(*model.LuckyGeneral))
	content += "\n✅抽奖活动已发布！\n"
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙返回到抽奖", "lucky$"),
			tgbotapi.NewInlineKeyboardButtonData("查看抽奖记录", "lucky_records"),
		),
	)
	// 1. create lucky activity
	rewards, _ := json.Marshal(data.Rewards)
	cond, _ := json.Marshal(map[string]interface{}{
		"users":     data.Users,
		"startTime": time.Now().Unix(),
		"endTime":   0,
	})
	item := model.LuckyActivity{
		ChatId:       chatId,
		LuckyName:    data.Name,
		LuckyType:    model.LuckyTypeGeneral,
		LuckySubType: data.SubType,
		LuckyCond:    string(cond),
		TotalReward:  "{}",
		Status:       model.LuckyStatusStart,
		RewardDetail: string(rewards), // 奖励信息 json
		StartTime:    time.Now().Unix(),
		EndTime:      0,
		PushChannel:  *data.Push,
	}
	services.CreateLucky(&item)
	// 2. push lucky info to chat group

	notify := tgbotapi.NewMessage(sess.groupChatId, buildLuckyMarkdown(bot, sess.groupChatId, sess.userId, data))
	if _, err := bot.Send(notify); err != nil {
		logger.Err(err).Stack().Msg("send lucky notify failed")
	}

	// 3. send reply
	reply := tgbotapi.NewEditMessageTextAndMarkup(sess.chatId, cb.Message.MessageID, content, inlineKeyboard)
	// reply.ReplyMarkup = inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Stack().Msg("send msg failed")
	}

	return nil
}

func buildLuckyMarkdown(bot *tgbotapi.BotAPI, chatId, userId int64, data *model.LuckyGeneral) string {
	/*
	   	🎁bigwinner 发起了通用抽奖活动

	   hhh
	   ├开奖时间：2023-09-07 01:01:00
	   ├参与关键词：andy
	   ├奖品列表：
	   ├       100USDT     ×10份

	   【如何参与？】在群组中回复关键词『andy』参与活动。
	*/
	var username string
	mgr := GroupManager{bot}
	user, err := mgr.fetchAndSaveMember(chatId, userId)
	if err != nil {
		username = fmt.Sprint(userId)
	} else {
		username = getDisplayName(&user)
	}
	content := "🎁" + mentionUser(username, userId) + " 发起了通用抽奖活动\n\n" + buildRewardInfo(data)
	content += fmt.Sprintf("\n【如何参与？】在群组中回复关键词『%s』参与活动。\n", data.Keyword)

	content = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, content)

	return content
}

func buildLuckyRecord(record *model.LuckyActivity) string {
	content := record.LuckyName + "\n"
	switch record.LuckyType {
	case model.LuckyTypeGeneral:
		content += fmt.Sprintf("├满人开奖  (%d人)\n")
		content += fmt.Sprintf("├参与关键词:  %s\n")
		content += fmt.Sprintf("├推送至频道:  %s\n")
		content += fmt.Sprintf("├奖品列表：\n")
	}

	content += fmt.Sprintf("\n创建者：%s\n", mentionUser(record.Creator, record.UserId))
	content += fmt.Sprintf("创建时间：%s\n", yyyymmddhhmmss(record.StartTime))
	content += fmt.Sprintf("状态: %s 已参与: %d人\n\n", record.Status, record.Participant)
	return content
}
