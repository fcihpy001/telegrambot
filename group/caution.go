package group

import (
	"errors"
	"fmt"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 触发警告及后续处理

var (
	chatBanWords = map[int64]*model.ProhibitedSetting{} // 群组违禁词
	chatBanNames = map[int64]*model.UserCheck{}         // 群组
	userCautions = map[string]int{}                     // key: cautionKey()
)

func cautionKey(chatId, userId int64, trigger model.TriggerType) string {
	return fmt.Sprintf("%d:%d:%v", chatId, userId, trigger)
}

func LoadChatRules() {
	prohibitedSettings, err := services.GetAllProhibitSettings()
	if err != nil {
		logger.Err(err).Msg("load all ProhibitSettings failed")
	} else {
		for _, item := range prohibitedSettings {
			chatBanWords[item.ChatId] = &item
		}
	}

	// 用户名
	banNames, err := services.GetAllUserCheck()
	if err != nil {
		logger.Err(err).Msg("load all userCheck failed")
	} else {
		for _, item := range banNames {
			chatBanNames[item.ChatId] = &item
		}
	}

	// 用户被警告次数
	cautions, err := services.GetAllCautions()
	if err != nil {
		logger.Err(err).Msg("load all cautions failed")
	} else {
		for _, item := range cautions {
			userCautions[cautionKey(item.ChatId, item.UserId, model.TriggerType(item.TriggerType))] = int(item.TriggerCount)
		}
	}
}

// 检查是否触发违禁词
func (mgr *GroupManager) onChatUserMessage(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	chatId := chat.ID
	wordRule := chatBanWords[chatId]
	if wordRule == nil || !wordRule.Enable {
		// 该群组没有违禁词过滤或未启用
		return
	}

	// todo 多个违禁词是怎么保存的
	words := strings.Split(wordRule.World, " ")
	text := msg.Text
	for _, word := range words {
		if !strings.Contains(text, word) {
			continue
		}
		// 触发违禁词 多个算一次
		if wordRule.Punish == model.PunishTypeWarning {
			mgr.doWarnPunish(chatId,
				msg.From.ID,
				model.TriggerTypeWords,
				wordRule.Punish,
				wordRule.WarningCount,
				wordRule.WarningAfterPunish,
				wordRule.BanTime,
				nil) // todo
		} else if wordRule.Punish == model.PunishTypeRevoke {
			mgr.deleteMessage(chatId, int64(msg.MessageID))
		} else {
			mgr.doPunish(chatId,
				msg.From.ID,
				model.TriggerTypeWords,
				wordRule.Punish,
				wordRule.BanTime,
				nil) // todo
		}

		return
	}
}

// 实施惩罚措施
func (mgr *GroupManager) doWarnPunish(
	chatId, userId int64,
	triggerType model.TriggerType,
	punishType model.PunishType,
	warnCount int,
	morePunishType model.PunishType,
	banTime int,
	punishArgs ...interface{},
) {
	key := cautionKey(chatId, userId, triggerType)
	count := userCautions[key]
	if count >= warnCount {
		// 1. 触发下一级惩罚
		// 2. 记录数清0
		count = 0
		mgr.doPunish(chatId, userId, triggerType, morePunishType, banTime, punishArgs...)
	} else {
		count++
		userCautions[key] = count
		// services.IncUserCautions()
	}
	if err := services.UpdateUserCaution(chatId, userId, triggerType, count); err != nil {
		logger.Err(err).Int64("chatId", chatId).Int64("userId", userId).Msg("update user caution failed")
	}
}

func (mgr *GroupManager) doPunish(
	chatId, userId int64,
	triggerType model.TriggerType,
	punishType model.PunishType,
	banTime int,
	punishArgs ...interface{},
) {
	ts := time.Now().Unix() + int64(banTime)
	switch punishType {
	case model.PunishTypeBan:
		// 禁言
		mgr.muteChatMember(chatId, userId, ts)

	case model.PunishTypeKick:
		// 踢出但不封禁
		ts := time.Now().Unix() + 40
		mgr.banChatMember(chatId, userId, ts)

	case model.PunishTypeBanAndKick:
		mgr.banChatMember(chatId, userId, ts)

	case model.PunishTypeRevoke:
		// 撤回消息

	default:
		logger.Err(errors.New("invalid punish type")).Msgf("punish type: %v", punishType)
	}
}

func (mgr *GroupManager) deleteMessage(chatId, msgId int64) error {
	msg := tgbotapi.DeleteMessageConfig{
		ChatID:    chatId,
		MessageID: int(msgId),
	}
	_, err := mgr.bot.Request(msg)
	if err != nil {
		logger.Err(err).Int64("chatId", chatId).Msg("delete message failed")
	}
	return err
}
