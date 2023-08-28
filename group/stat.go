package group

import (
	"encoding/json"
	"telegramBot/model"
	"telegramBot/services"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		return
	}
	rows, err := services.FindChatMessageCount(model.StatTypeMessageCount, chat.ID, startTs, endTs, 0, 5)
	if err != nil {
		logger.Err(err)
		return
	}
	_ = rows
}

// just for test
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
		Name:               "xxxxxxxx",
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
