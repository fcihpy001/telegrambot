package group

import (
	"encoding/json"
	"fmt"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DoStat 统计入口
//
//	消息统计
//	进群统计
//	邀请统计
//	离群统计
func DoStat(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		msg := update.Message
		if msg.IsCommand() {
			return
		}
		if msg.From != nil {
			if msg.NewChatMembers != nil {
				// new group member
				services.StatsNewMembers(update)
				for _, member := range msg.NewChatMembers {
					if member.ID == utils.GetBotUserId() {
						// 第一次被邀请进入群, 获取群信息及群用户
						mgr := GroupManager{bot}
						mgr.GetChatInfo(msg.Chat.ID)
					}
				}
				return
			}
			if msg.LeftChatMember != nil {
				// leave group
				services.StatsLeave(update)
				return
			}

			// message stat
			if msg.From.IsBot {
				return
			}
			chat := msg.Chat
			if chat != nil && !chat.IsPrivate() {
				// 消息统计
				services.StatChatMessage(chat.ID, msg.From.ID, int64(msg.Date))
				return
			}
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

// 获取chat 群用户数
func (mgr *GroupManager) GetChatMemberCount(id int64) int {
	resp, err := mgr.bot.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{ChatConfig: tgbotapi.ChatConfig{
		ChatID: id,
	}})
	if err != nil {
		logger.Err(err).Int64("chatId", id).Msg("get chat member count failed")
		return 0
	}

	return resp
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
	rows, err := services.FindChatCountGroupByUser(model.StatTypeMessageCount, chat.ID, startTs/60, endTs/60, 0, 5)
	// rows, err := services.FindChatCountGroupByUser(model.StatTypeMessageCount, chat.ID, startTs, endTs, 0, 5)
	if err != nil {
		logger.Err(err)
		return
	}
	res := ""
	for _, row := range rows {
		res += fmt.Sprintf("%d    %d\n", row.UserId, row.Count)
	}
	mgr.sendText(chat.ID, res)
}

//lint:ignore U1000 just for test
func (mgr *GroupManager) inviteLink(update *tgbotapi.Update) {
	msg := update.Message
	if msg.Chat == nil {
		logger.Warn().Msg("not chat group")
		return
	}
	chatId := msg.Chat.ID
	resp := tgbotapi.CreateChatInviteLinkConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatId,
		},
		Name:               "fc",
		ExpireDate:         int(time.Now().Unix() + 86400*365),
		MemberLimit:        9999,
		CreatesJoinRequest: false,
	}
	link, err := mgr.bot.Request(resp)
	if err != nil {
		logger.Warn().Msgf("invite send failed: %v", err)
	}

	m := map[string]interface{}{}
	json.Unmarshal(link.Result, &m)
	// fmt.Println(prettyJSON(link))
	inviteMsg := tgbotapi.NewMessage(chatId, m["invite_link"].(string))
	mgr.sendMessage(inviteMsg, "send invite link failed")
}
