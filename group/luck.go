package group

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
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

var (
	luckyEndChan  chan int
	luckyCreated  chan *model.LuckyActivity
	luckyLock     sync.RWMutex
	luckyKeywords = map[string][]*model.LuckyActivity{}

	_bot *tgbotapi.BotAPI
)

func SetBot(botapi *tgbotapi.BotAPI) {
	_bot = botapi
}

// 监听所有 lucky keywords
func InitLuckyFilter(ctx context.Context) {
	luckies := services.GetAllLuckyActivities()

	for _, item := range luckies {
		luckyKeywords[item.Keyword] = append(luckyKeywords[item.Keyword], item)
	}

	luckyEndChan = make(chan int, 1)
	luckyCreated = make(chan *model.LuckyActivity, 1)

	tmr := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("context cancel")
				return
			case <-tmr.C:
				loopLuckyKeywords()
			case <-luckyEndChan:
				loopLuckyKeywords()
			case item := <-luckyCreated:
				luckyKeywords[item.Keyword] = append(luckyKeywords[item.Keyword], item)
				logger.Info().Str("luckyName", item.LuckyName).Msg("lucky created")
			}
		}
	}()
}

func loopLuckyKeywords() {
	luckyLock.Lock()
	defer luckyLock.Unlock()

	now := time.Now().Unix()
	for word, records := range luckyKeywords {
		nRecords := []*model.LuckyActivity{}
		for _, record := range records {
			if record.LuckySubType == model.LuckySubTypeTime && record.EndTime < now {
				// record is time up
				record.Status = model.LuckyStatusEnd
				// 这里需要 bot 实例
				luckyOpenReward(_bot, record)
			}
			if record.Status == model.LuckyStatusStart {
				nRecords = append(nRecords, record)
			}
		}
		luckyKeywords[word] = nRecords
		if len(nRecords) == 0 {
			delete(luckyKeywords, word)
		}
	}
}

// 开奖
func luckyOpenReward(bot *tgbotapi.BotAPI, record *model.LuckyActivity) {
	var rewards []model.LuckyReward
	shares := 0
	json.Unmarshal([]byte(record.RewardDetail), &rewards)
	for _, reward := range rewards {
		shares += reward.Shares
	}
	flatRewards := make([]model.LuckyReward, shares)
	idx := 0
	used := 0
	for i := 0; i < shares; i++ {
		flatRewards[i] = rewards[idx]
		used++
		if used >= rewards[idx].Shares {
			used = 0
			idx++
		}
	}

	parts := services.GetLuckyAllParticipates(record)
	luckies := []model.LuckyRecord{} // 中奖用户
	if len(parts) > 0 {
		counter := len(parts)
		rewardIdx := 0
		for i := 0; i < len(parts); {
			val := rand.Intn(counter)
			if rewardIdx >= len(flatRewards) {
				// 奖金发完
				break
			}
			if parts[val].Reward != "" {
				// 已经中奖
				continue
			} else {
				parts[val].Reward = flatRewards[rewardIdx].Name
				rewardIdx++
				i++
			}
		}
		// 更新数据库
		rewardParts := 0
		for _, part := range parts {
			if part.Reward != "" {
				services.UpdateLuckyRewardRecord(&part)
				rewardParts++
				luckies = append(luckies, part)
			}
		}
		record.PartReward = rewardParts
		record.RewardRatio = fmt.Sprint(len(flatRewards)*100/rewardParts) + "%"
	}
	services.UpdateLuckyActivity(record)

	/* 中奖结果通知
		🎁活动暴富go 开奖啦！
	总共参与2人，综合中奖率50%

	🥳🥳恭喜以下中奖用户：

	🎉bigwinner 获得奖品：100usdt

	👮🏼 抽奖创建者：bigwinner
	『联系该群管理领取您的奖品吧~』
	🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉
	*/
	content := escapeText(fmt.Sprintf("🎁活动[%s] 开奖啦！\n总共参与%d人，综合中奖率%s\n\n",
		record.LuckyName, record.Participant, record.RewardRatio))
	content += "🥳🥳恭喜以下中奖用户：\n\n"
	for _, item := range luckies {
		content += "🎉" + mentionUser(item.Username, item.UserId) + " 获得奖品：" + escapeText(item.Reward) + "\n"
	}
	content += "\n👮🏼 抽奖创建者：" + mentionUser(record.Creator, record.UserId) + "\n"
	content += escapeText("『联系该群管理领取您的奖品吧~』\n🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉\n")
	sendMarkdown(bot, record.ChatId, content, true)
}

// 记录数据库
// 判断抽奖是否达到结束条件
func onLuckyTrigger(update *tgbotapi.Update, bot *tgbotapi.BotAPI, record *model.LuckyActivity) {
	logger.Info().Msg("luck triggered")

	msg := update.Message
	fromId := msg.From.ID
	// 用户是否已经参与过
	if services.CheckUserHasParticipated(int64(record.ID), fromId) {
		//
		reject := tgbotapi.NewMessage(msg.Chat.ID, "您已参加过该活动，请勿重复参加！")
		reject.ReplyToMessageID = msg.MessageID
		resp, err := bot.Send(reject)
		if err != nil {
			logger.Err(err).Msg("send message failed")
		} else {
			// delete message
			setTimer(30, func() {
				sendDeleteMsg(bot, msg.Chat.ID, resp.MessageID)
			})
		}
		return
	}

	go services.OnLuckyParticipate(record, fromId, getDisplayNameFromUser(update.Message.From))

	record.Participant += 1
	// 发送参与通知
	reply := tgbotapi.NewMessage(msg.Chat.ID,
		buildParticiateContent(record, update))
	reply.ReplyToMessageID = msg.MessageID
	reply.ParseMode = tgbotapi.ModeMarkdownV2
	resp, err := bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send participate lucky notify failed")
	} else {
		// delete message
		setTimer(30, func() {
			sendDeleteMsg(bot, msg.Chat.ID, resp.MessageID)
		})
	}

	if record.ReachParticipantUsers() {
		logger.Info().Uint("luckyId", record.ID).Msgf("lucky [%s] reach users", record.LuckyName)
		record.Status = model.LuckyStatusEnd
		go luckyOpenReward(bot, record)
	}
}

func MatchLuckyKeywords(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}
	text := update.Message.Text

	changed := false
	luckyLock.RLock()
	for word, records := range luckyKeywords {
		if text == word {
			// trigger record
			for _, record := range records {
				onLuckyTrigger(update, bot, record)
				if record.Status != model.LuckyStatusStart {
					changed = true
				}
			}
		}
	}
	luckyLock.RUnlock()

	if changed {
		luckyEndChan <- 1
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
				tgbotapi.NewInlineKeyboardButtonData("⬅️上一条", fmt.Sprintf("lucky_record?idx=%d", idx-1)),
				tgbotapi.NewInlineKeyboardButtonData("下一条➡️", fmt.Sprintf("lucky_record?idx=%d", idx+1)),
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
	reply.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send lucky record failed")
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
	//
	total, opened, canceled := services.StatChatLuckyCount(param.chatId)
	msg := tgbotapi.NewEditMessageTextAndMarkup(param.chatId, param.msgId,
		fmt.Sprintf("🎁【测试】抽奖\n\n发起抽奖次数：%d    \n\n已开奖：%d       未开奖：%d       取消：%d",
			total, opened, total-opened-canceled, canceled),
		inlineKeyboard)

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
				tgbotapi.NewInlineKeyboardButtonData("满人开奖", "lucky_create_general?subType=users"),
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

func luckyCancel(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error {
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
	record := services.GetLuckyActivity(chatId, int(idx))
	if record.Status != model.LuckyStatusStart {
		logger.Warn().Str("status", record.Status).Str("idx", sidx).Msg("cannot cancel")
		return nil
	}
	record.Status = model.LuckyStatusCancel
	services.UpdateLuckyActivity(record)
	luckyRecords(update, bot, param)
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
	case model.LuckySubTypeUsers:
		// 满人抽奖
		content = "🎁创建通用抽奖(/cancel 命令返回首页)\n\n" +
			"请回复参与多少人后开奖：\n\n"
	case model.LuckySubTypeTime:
		// 定时抽奖
		content = "🎁创建通用抽奖(/cancel 命令返回首页)\n\n" +
			"请回复开奖时间：\n" +
			"格式：年-月-日 时:分\n" +
			"例如：2023-09-11 19:45\n\n"
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
	// StartAdminConversation(param.chatId, param.chatId, update.CallbackQuery.From.ID, int64(param.msgId),
	// 	ConversationLuckyCreateGeneralStep1,
	// 	&data,
	// 	luckyCreateGeneralSteps,
	// )
	updateAdminConversation(param.chatId,
		ConversationLuckyCreateGeneralStep1,
		&data,
		luckyCreateGeneralSteps)

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
	return escapeText("🎁创建通用抽奖  ( /cancel 命令返回首页)\n\n") + buildRewardBody(data)
}

func buildRewardBody(data *model.LuckyGeneral) string {
	content := ""
	if data.Name != "" {
		content += escapeText(data.Name) + "\n"
	}
	if data.Keyword != "" {
		content += fmt.Sprintf("├参与关键词：%s\n", escapeText(data.Keyword))
	}
	if data.Push != nil {
		if *data.Push {
			content += "├推送至频道：❌\n"
		} else {
			content += "├推送至频道：✅\n"
		}
	}
	if data.SubType == model.LuckySubTypeUsers {
		content += escapeText(fmt.Sprintf("├满人开奖  (%v人)\n", data.Users))
	} else if data.SubType == model.LuckySubTypeTime {
		content += escapeText(fmt.Sprintf("├开奖时间:  (%v)\n", yyyymmddhhmmss(data.EndTime)))
	}
	content += "├奖品列表:\n"
	for _, reward := range data.Rewards {
		if reward.Shares > 0 {
			content += fmt.Sprintf("├       %s    x %d份\n", escapeText(reward.Name), reward.Shares)
		} else {
			content += fmt.Sprintf("├       %s\n", escapeText(reward.Name))
		}
	}

	return content
}

func buildParticiateContent(record *model.LuckyActivity, update *tgbotapi.Update) string {
	content := "🎁" + escapeText(record.LuckyName) + "\n\n"
	msg := update.Message

	username := getDisplayNameFromUser(msg.From)
	content += mentionUser(username, msg.From.ID) + " 您已参与成功，请等待开奖通知！\n\n"

	if record.LuckyType == model.LuckyTypeGeneral && record.LuckySubType == model.LuckySubTypeUsers {
		content += escapeText(fmt.Sprintf("├%s  \\(%d人\\)\n", record.GetLuckyType(), record.GetLuckGeneralUsers()))
	} else {
		content += escapeText(fmt.Sprintf("├开奖时间:  \\(%s\\)\n", yyyymmddhhmmss(record.EndTime)))
	}
	content += fmt.Sprintf("├已参与  \\(%d人\\)\n", record.Participant)
	content += fmt.Sprintf("├参与关键词：  %s\n", escapeText(record.Keyword))
	content += "├奖品列表：\n"
	for _, reward := range record.GetRewards() {
		content += fmt.Sprintf("├    %s  x %d份\n", escapeText(reward.Name), reward.Shares)
	}

	content += fmt.Sprintf("\n【如何抽奖？】在群组中回复关键词『%s』参与活动。\n", escapeText(record.Keyword))
	return content
}

// 用于展示抽奖活动
func buildRewardInfo(data *model.LuckyGeneral) string {
	content := fmt.Sprintf("%s\n├开始时间：%s\n├参与关键词：%s\n├奖品列表：\n",
		escapeText(data.Name),
		escapeText(yyyymmddhhmmss(data.StartTime)),
		escapeText(data.Keyword),
	)
	for _, reward := range data.Rewards {
		if reward.Shares > 0 {
			content += fmt.Sprintf("├       %s    x %d份\n", escapeText(reward.Name), reward.Shares)
		} else {
			content += fmt.Sprintf("├       %s\n", escapeText(reward.Name))
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
		var content string
		sess.status = ConversationLuckyCreateGeneralStep2
		if data.SubType == model.LuckySubTypeUsers {
			users, err := strconv.Atoi(text)
			if err != nil {
				logger.Err(err).Msg("invalid user arg")
			}
			data.Users = users
			content = escapeText(fmt.Sprintf("🎁创建通用抽奖  ( /cancel 命令返回首页)\n\n├满人开奖  (%s人)\n\n请回复第一个奖品的名称（如：1USDT）：", text))
		} else {
			tm, err := parseDateTime(text)
			if err != nil {
				logger.Err(err).Msg("invalid lucky end time")
			}
			if tm.Unix() <= time.Now().Unix() {
				logger.Error().Msg("lucky end time less than current time")
			}
			data.EndTime = tm.Unix()
			content = escapeText(fmt.Sprintf("🎁创建通用抽奖  ( /cancel 命令返回首页)\n\n├开奖时间:  (%s)\n\n请回复第一个奖品的名称（如：1USDT）：", text))
		}
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
		"endTime":   data.EndTime,
	})
	data.StartTime = time.Now().Unix()
	item := model.LuckyActivity{
		ChatId:       chatId,
		LuckyName:    data.Name,
		LuckyType:    model.LuckyTypeGeneral,
		LuckySubType: data.SubType,
		UserId:       cb.Message.From.ID,
		Creator:      getDisplayNameFromUser(cb.Message.From),
		Keyword:      data.Keyword,
		LuckyCond:    string(cond),
		TotalReward:  "{}",
		Status:       model.LuckyStatusStart,
		RewardDetail: string(rewards), // 奖励信息 json
		StartTime:    time.Now().Unix(),
		EndTime:      data.EndTime,
		PushChannel:  *data.Push,
	}
	services.CreateLucky(&item)
	luckyCreated <- &item
	// 2. push lucky info to chat group
	notifyText := buildLuckyMarkdown(bot, sess.groupChatId, sess.userId, data)
	sendMarkdown(bot, sess.groupChatId, notifyText, true)

	// 3. send reply
	reply := tgbotapi.NewEditMessageTextAndMarkup(sess.chatId, cb.Message.MessageID, content, inlineKeyboard)
	// reply.ReplyMarkup = inlineKeyboard
	reply.ParseMode = tgbotapi.ModeMarkdownV2
	if _, err := bot.Send(reply); err != nil {
		logger.Err(err).Stack().Str("content", content).Msg("send msg failed")
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
	content += fmt.Sprintf("\n【如何参与？】在群组中回复关键词『%s』参与活动。\n", escapeText(data.Keyword))

	return content
}

func buildLuckyRecord(record *model.LuckyActivity) string {
	content := escapeText(record.LuckyName + "\n")
	switch record.LuckyType {
	case model.LuckyTypeGeneral:
		var (
			cond    map[string]interface{}
			rewards []model.LuckyReward
		)
		json.Unmarshal([]byte(record.LuckyCond), &cond)
		json.Unmarshal([]byte(record.RewardDetail), &rewards)
		content += escapeText(fmt.Sprintf("├满人开奖  (%d人)\n", int(cond["users"].(float64))))
		content += fmt.Sprintf("├参与关键词:  %s\n", escapeText(record.Keyword))
		content += fmt.Sprintf("├推送至频道:  %s\n", "❌")
		content += "├奖品列表：\n"
		for _, reward := range rewards {
			content += fmt.Sprintf("├       %s    x %d份\n", escapeText(reward.Name), reward.Shares)
		}
	}

	content += fmt.Sprintf("\n创建者：%s\n", mentionUser(record.Creator, record.UserId))
	content += fmt.Sprintf("创建时间：%s\n", escapeText(yyyymmddhhmmss(record.StartTime)))
	content += fmt.Sprintf("状态: %s 已参与: %d人\n\n", luckyStatus(record.Status), record.Participant)
	return content
}

func luckyStatus(status string) string {
	switch status {
	case model.LuckyStatusStart:
		return "进行中✅"
	case model.LuckyStatusCancel:
		return "已取消❌"
	case model.LuckyStatusEnd:
		return "已开奖 ⭕️"
	}

	return status
}

func LuckyCreateCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("欢迎使用@%s：\n\n点击下面按钮创建抽奖(仅限管理员)", bot.Self.UserName)
	url := fmt.Sprintf("https://t.me/%s?start=lucky_%d", bot.Self.UserName, update.Message.Chat.ID)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("👉🎁 点击创建抽奖活动👈", url)))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
