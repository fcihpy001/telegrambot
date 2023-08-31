package group

import (
	"fmt"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mgr *GroupManager) statics(update *tgbotapi.Update) {
	btn11 := model.ButtonInfo{
		Text:    ".发言统计.",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn21 := model.ButtonInfo{
		Text:    "今日发言排名",
		Data:    "group_speechtodayranging",
		BtnType: model.BtnTypeData,
	}
	btn22 := model.ButtonInfo{
		Text:    "7日发言排名",
		Data:    "group_speech7daysranging",
		BtnType: model.BtnTypeData,
	}
	btn23 := model.ButtonInfo{
		Text:    "7日发言统计",
		Data:    "group_speechstatistics",
		BtnType: model.BtnTypeData,
	}
	btn31 := model.ButtonInfo{
		Text:    "📊邀请统计",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn41 := model.ButtonInfo{
		Text:    "今日邀请排名",
		Data:    "group_invite_ranging",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "7日邀请排名",
		Data:    "group_invite_7days_ranging",
		BtnType: model.BtnTypeData,
	}
	btn61 := model.ButtonInfo{
		Text:    "📊进退群统计",
		Data:    "toast",
		BtnType: model.BtnTypeData,
	}
	btn71 := model.ButtonInfo{
		Text:    "今日进退群数据",
		Data:    "group_today_quit",
		BtnType: model.BtnTypeData,
	}
	btn72 := model.ButtonInfo{
		Text:    "7日进退群统计",
		Data:    "group_7days_quit",
		BtnType: model.BtnTypeData,
	}
	btn81 := model.ButtonInfo{
		Text:    "返回首页",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	btns := [][]model.ButtonInfo{{btn11}, {btn21, btn22, btn23}, {btn31}, {btn41, btn42}, {btn61}, {btn71, btn72}, {btn81}}
	keyboard := utils.MakeKeyboard(btns)
	utils.StaticsMarkup = keyboard
	content := "📊 【流量聚集地】统计\n\n在群组中使用命令：\n/stat 查询今天活跃统计\n/stat_week 查询七天活跃统计\n/stats 自定义时间查询活跃统计\n\n查看命令帮助"
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	mgr.bot.Send(msg)
}

func (mgr *GroupManager) speechRanging(update *tgbotapi.Update, period string) {
	startTs, endTs := getTimeRange(period)
	chatId := update.CallbackQuery.Message.Chat.ID

	result, err := mgr.statPeroidChatMessages(chatId, startTs, endTs, 0, 10)
	if err != nil {
		logger.Err(err).Msg("speechRanging failed")
		return
	}
	content := ""
	fmtUserRating(1, result.Data)

	if period == "week" {
		content = fmt.Sprintf("7日发言数：%d 条 聊天人数: %d，以下是排名: \n\n", result.TotalMsg, result.TotalUser)
	} else {
		content = fmt.Sprintf("今日总发言：%d 条，聊天人数: %d, 以下是排名：\n\n", result.TotalMsg, result.TotalUser)
	}
	content += fmtUserRating(1, result.Data)
	mgr.staticsDetail(update, content)
}

func (mgr *GroupManager) speechstatistics(update *tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID
	startDay, endDay := getWeekRange()
	records, err := services.GroupChatMessageByDay(chatId, startDay, endDay)
	if err != nil {
		logger.Err(err).Msg("stat messages by day failed")
		return
	}
	content := "以下是7日发言统计:\n"
	for _, record := range records {
		content += fmt.Sprintf("%s   %d\n", record.Day, record.Count)
	}
	mgr.staticsDetail(update, content)
}

// 今日邀请
func (mgr *GroupManager) invitesToday(update *tgbotapi.Update) {
	startTs, endTs := getTimeRange("today")
	mgr.invitestatis(update, startTs, endTs, 0, "今日")
}

func (mgr *GroupManager) invitestatis(update *tgbotapi.Update, startTs, endTs int64, startIdx int, timeRange string) {
	chatId := update.CallbackQuery.Message.Chat.ID
	total, invites, err := services.GroupChatInviteByUser(chatId, startTs, endTs, 10, int64(startIdx))
	if err != nil {
		logger.Err(err).Msg("stat today invites failed")
		return
	}
	var ids []int64
	for _, record := range invites {
		ids = append(ids, record.InvitedBy)
	}
	names := mgr.getUserNames(chatId, ids)
	content := fmt.Sprintf("%s共邀请: %d人，以下是排名\n", timeRange, total)
	for i, record := range invites {
		name := names[record.InvitedBy]
		if name == "" {
			name = fmt.Sprintf("%d", record.InvitedBy)
		}
		content += fmt.Sprintf("%d\\. \t%s \\-  %d\n", 1+startIdx+i, mentionUser(name, record.InvitedBy), record.Count)
	}
	// 1.fcihpy - 6 条\n2.Fcihpy3 - 1 条"
	mgr.staticsDetail(update, content)
}

// 7日邀请
func (mgr *GroupManager) invitesWeek(update *tgbotapi.Update) {
	startTs, endTs := getTimeRange("week")
	mgr.invitestatis(update, startTs, endTs, 0, "7日")
}

// 今日进群数据
func (mgr *GroupManager) groupmemberstatis(update *tgbotapi.Update, period string) {
	chatId := update.CallbackQuery.Message.Chat.ID
	startTs, endTs := getTimeRange(period)
	var content string
	if period == "today" {
		joinCount, _ := services.CountChatJoinLeft("join", chatId, startTs, endTs)
		leftCount, _ := services.CountChatJoinLeft("left", chatId, startTs, endTs)
		joinList, leftList, err := services.GetLatestJoinLeftUsers(chatId, startTs, endTs, 10)
		if err != nil {
			logger.Err(err).Msg("get latest join/left users failed")
			return
		}
		content = fmt.Sprintf("今日进群: %d 人，退群: %d人\n", joinCount, leftCount)
		content += "\n以下是今日最新进群用户:\n"
		// 以下是今日最新进群20人：\n\n\n以下是今日最新退群20人：
		for _, item := range joinList {
			content += mentionUser(getDisplayName(&item), item.Uid) + "\n"
		}
		content += "\n以下是今日最新退群用户:\n"
		for _, item := range leftList {
			content += mentionUser(getDisplayName(&item), item.Uid) + "\n"
		}
	} else {
		content = "7日进群：%d人，退群：%d人\n以下是近7日进群退群人数:\n日期 \t  进群数量    退群数量\n"
		joins, _ := services.GroupChatJoinLeftByDay("join", chatId, startTs, endTs)
		leaves, _ := services.GroupChatJoinLeftByDay("left", chatId, startTs, endTs)
		days := getRangeDays(startTs, endTs)
		joinIdx := 0
		leaveIdx := 0
		for _, day := range days {
			join := 0
			left := 0
			if len(joins) > joinIdx && joins[joinIdx].Day == day {
				join = joins[joinIdx].Count
				joinIdx++
			}
			if len(leaves) > leaveIdx && leaves[leaveIdx].Day == day {
				left = leaves[leaveIdx].Count
				leaveIdx++
			}
			content += fmt.Sprintf("%s   %d   %d\n", day, join, left)
		}
	}

	mgr.staticsDetail(update, content)
}

func (mgr *GroupManager) staticsDetail(update *tgbotapi.Update, content string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦀返回", "group_back_statics"),
		))
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	msg.ParseMode = "MarkdownV2"
	_, err := mgr.bot.Send(msg)
	if err != nil {
		logger.Err(err).Stack().Str("content", content).Msg("send msg failed")
	}
}

func (mgr *GroupManager) group_back_statics(update *tgbotapi.Update) {
	txt := "统计数据"
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, txt, utils.StaticsMarkup)
	mgr.bot.Send(msg)
}
