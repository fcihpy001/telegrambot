package group

import (
	"fmt"
	"telegramBot/model"
	"telegramBot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserCounter struct {
	UserId   int64
	Count    int64
	UserName string
}

type DayCounter struct {
	Day   string
	Count int64
}

type StatMsgResult struct {
	StartTs   int64
	EndTs     int64
	TotalMsg  int // 发言总数
	TotalUser int // 发言过的用户总数
	Data      []UserCounter
}

// DoStat 统计入口
//
//	消息统计
//	进群统计
//	邀请统计
//	离群统计
func DoStat(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		msg := update.Message
		if msg.IsCommand() || msg.From == nil {
			return
		}
		// 进群统计
		if msg.NewChatMembers != nil {
			services.StatsNewMembers(update)
			for _, member := range msg.NewChatMembers {
				if member.ID == bot.Self.ID {
					// 第一次被邀请进入群, 获取群信息及群用户
					mgr := GroupManager{bot}
					mgr.GetChatInfo(msg.Chat.ID)
				}
			}
			return
		}
		// 离群统计
		if msg.LeftChatMember != nil {
			services.StatsLeave(update)
			return
		}

		// 消息统计
		if msg.From.IsBot {
			return
		}
		chat := msg.Chat
		if chat != nil && !chat.IsPrivate() {
			services.StatChatMessage(chat.ID, msg.From.ID, int64(msg.Date))
			return
		}
	}
}

// 获取群及群用户信息
func (mgr *GroupManager) GetChatInfo(id int64) {
	resp, err := mgr.bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: id,
		},
	})

	if err != nil {
		logger.Err(err).Int64("chatId", id).Msg("get chat info failed")
		return
	}
	fmt.Println(prettyJSON(resp))
	// save chat group
	services.SaveChatGroupByChat(&resp)
}

func (mgr *GroupManager) statPeroidChatMessages(chatId, startTs, endTs int64, offset, limit int) (StatMsgResult, error) {
	totalMsg, totalUser := services.CountTotalChatMessageUsers(chatId, startTs, endTs)
	result := StatMsgResult{
		StartTs:   startTs,
		EndTs:     endTs,
		TotalMsg:  totalMsg,
		TotalUser: totalUser,
	}
	rows, err := services.GetMessageCountGroupByUser(chatId, startTs, endTs, int64(offset), int64(limit))
	if err != nil {
		return result, err
	}
	result.Data = mgr.convertToUserCounter(chatId, rows)
	return result, nil
}

// 统计回应
func (mgr *GroupManager) StatsMemberMessages(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	if chat == nil {
		logger.Warn().Msg("not group chat message")
		return
	}
	startTs, endTs, err := parseTimeRange(msg.Text)
	if err != nil {
		mgr.sendText(chat.ID, err.Error())
		logger.Warn().Err(err).Msg("invalid time range")
		return
	}
	result, err := mgr.statPeroidChatMessages(chat.ID, startTs, endTs, 0, 5)
	if err != nil {
		logger.Err(err)
		return
	}
	res := fmtUserRating(1, result.Data)
	mgr.sendText(chat.ID, res)
}

// 邀请排名
func (mgr *GroupManager) StatInvite(update *tgbotapi.Update, startTs, endTs int64) {

}

// 进群、退群排名
func (mgr *GroupManager) StatJoinLeave(update *tgbotapi.Update, startTs, endTs int64) {

}

// 1 调整格式
// 2 用户名可以点击
func (mgr *GroupManager) convertToUserCounter(chatId int64, items []model.StatCount) (data []UserCounter) {
	var ids []int64
	for _, item := range items {
		ids = append(ids, item.UserId)
	}
	names := mgr.getUserNames(chatId, ids)

	fmt.Println("names:", ids, names)
	for _, item := range items {
		data = append(data, UserCounter{
			UserId:   item.UserId,
			Count:    int64(item.Count),
			UserName: names[item.UserId],
		})
	}
	return
}

func fmtUserRating(startIdx int, items []UserCounter) (res string) {
	for i, item := range items {
		if item.UserName != "" {
			res += fmt.Sprintf("%d\\. %s   %d\n", startIdx+i, mentionUser(item.UserName, item.UserId), item.Count)
		} else {
			res += fmt.Sprintf("%d\\. %s   %d\n", startIdx+i, mentionUser(item.UserId, item.UserId), item.Count)
		}
	}
	return
}

// 根据用户 id 查找用户名
// 首先从数据库中查找, 如果查找不到, 使用 getChatMember 接口查找并入库
func (mgr *GroupManager) getUserNames(chatId int64, ids []int64) map[int64]string {
	users := services.GetUserNames(ids)
	names := map[int64]string{}
	for _, user := range users {
		names[user.Uid] = getDisplayName(&user)
	}
	for _, id := range ids {
		if names[id] == "" {
			//
			user, err := mgr.fetchAndSaveMember(chatId, id)
			if err == nil {
				names[id] = getDisplayName(&user)
			}
		}
	}
	return names
}

func (mgr *GroupManager) getUserName(chatId int64, id int64) string {
	users := services.GetUserNames([]int64{id})
	if len(users) > 0 {
		user := users[id]
		return getDisplayName(&user)
	}
	user, err := mgr.fetchAndSaveMember(chatId, id)
	if err != nil {
		logger.Err(err).Int64("userId", id).Msg("fetch user info failed")
		return fmt.Sprint(id)
	}
	return getDisplayName(&user)
}

func getDisplayName(u *model.User) string {
	return u.FirstName + " " + u.LastName
}
