package group

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// 记录那个聊天会话的状态以及后续动作

type ConversationStatus string

const (
	ConversationStart              ConversationStatus = "start"
	ConversationWaitSolitaireInput ConversationStatus = "waitSolitaireInput"
)

type botAdminSession struct {
	groupChatId int64 // supergroup chat ID
	chatId      int64 // private conversation chat ID
	userId      int64
	status      ConversationStatus
}

var (
	adminSessions = map[int64]*botAdminSession{}
)

func startAdminConversation(groupChatId, chatId, userId int64) {
	adminSessions[groupChatId] = &botAdminSession{
		groupChatId: groupChatId,
		chatId:      chatId,
		userId:      userId,
		status:      ConversationStart,
	}
}

func updateAdminConversation(chatId int64, status ConversationStatus) {
	sess := adminSessions[chatId]
	if sess == nil {
		logger.Error().Int64("chatId", chatId).Msg("not found admin conversation")
		return
	}
	sess.status = status
}

// 最后的输入
func handleAdminConversation(update *tgbotapi.Update) {

}
