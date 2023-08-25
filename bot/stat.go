package bot

import (
	"encoding/json"
	"log"
	"telegramBot/model"
	"telegramBot/services"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *SmartBot) StatsMemberMessages(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	if chat == nil {
		log.Println("not group chat message")
		return
	}
	startTs, endTs, err := parseTimeRange(msg.Text)
	if err != nil {
		bot.SendText(chat.ID, err.Error())
		return
	}
	rows, err := services.FindChatMessageCount(model.StatTypeMessageCount, chat.ID, startTs, endTs, 0, 5)
	if err != nil {
		log.Println(err)
		return
	}
	_ = rows
}

// just for test
func (bot *SmartBot) inviteLink(update *tgbotapi.Update) {
	msg := update.Message
	if msg.Chat == nil {
		log.Printf("not chat group")
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
	link, err := bot.bot.Request(resp)
	if err != nil {
		log.Printf("invite send failed: %v", err)
	}

	m := map[string]interface{}{}
	json.Unmarshal(link.Result, &m)
	// fmt.Println(prettyJSON(link))
	inviteMsg := tgbotapi.NewMessage(chatId, m["invite_link"].(string))
	bot.sendMessage(inviteMsg, "send invite link failed")
}
