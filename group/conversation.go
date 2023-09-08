package group

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 记录那个聊天会话的状态以及后续动作

type ConversationStatus string

const (
	ConversationStart              ConversationStatus = "start"
	ConversationWaitSolitaireDesc  ConversationStatus = "waitSolitaireDesc"
	ConversationWaitSolitaireUsers ConversationStatus = "waitSolitaireUsers" // limitUser
	ConversationPlaySolitaire      ConversationStatus = "playSolitaire"      // limitUser
)

type ConversationFn func(update *tgbotapi.Update, bot *tgbotapi.BotAPI, sess *botConversation) error

type botConversation struct {
	groupChatId int64 // supergroup chat ID
	chatId      int64 // private conversation chat ID
	userId      int64
	messageId   int64
	status      ConversationStatus
	data        interface{}
	fn          ConversationFn
}

var sessions = map[int64]*botConversation{}

func GetConversation(chatId int64) *botConversation {
	return sessions[chatId]
}

func StartAdminConversation(groupChatId, chatId, userId, messageId int64,
	status ConversationStatus,
	data interface{},
	fn ConversationFn) {
	sessions[groupChatId] = &botConversation{
		groupChatId: groupChatId,
		chatId:      chatId,
		userId:      userId,
		messageId:   messageId,
		status:      status,
		data:        data,
		fn:          fn,
	}
}

func UpdateAdminConversation(chatId int64, status ConversationStatus) {
	sess := sessions[chatId]
	if sess == nil {
		logger.Error().Int64("chatId", chatId).Msg("not found admin conversation")
		return
	}
	sess.status = status
}

func RemoveConversation(chatId int64) {
	delete(sessions, chatId)
}

// 最后的输入
func HandleAdminConversation(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	// println("handleAdminConversation:" + prettyJSON(update))
	msg := update.Message
	userId := msg.From.ID
	chat := msg.Chat
	chatId := chat.ID
	conversion, ok := sessions[chatId]
	if !ok {
		//
		return false
	}
	if conversion.userId != userId {
		return false
	}

	// conversion.fn(update, bot, conversion)

	mgr := GroupManager{bot}
	switch conversion.status {
	case ConversationWaitSolitaireDesc:
		logger.Info().Int64("chatId", chatId).Msg("solitaire create completed")
		mgr.onSolitaireCreated(update, conversion)
		RemoveConversation(chatId)
		return true

	case ConversationWaitSolitaireUsers:
		logger.Info().Int64("chatId", chatId).Msg("solitaire limit users")
		mgr.onSolitaireLimitUser(update, conversion)
		// RemoveConversation(chatId)
		return true

	case ConversationPlaySolitaire:
		logger.Info().Int64("chatId", chatId).Msg("play solitaire message")
		mgr.onPlaySolitaireComplete(update, conversion)
		RemoveConversation(chatId)
		return true
	}
	return false
}
