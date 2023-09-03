package group

import (
	"fmt"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 接龙

var (
	solitaireStatus = map[string]string{
		model.SolitaireStatusActive: "收集中",
		model.SolitaireStatusEnded:  "已结束",
	}
)

// 接龙首屏
func (mgr *GroupManager) SolitaireIndex(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	chatId := chat.ID
	items, err := services.GetChatSolitaireList(chatId)
	if err != nil {
		logger.Err(err).Msg("get solitaire list failed")
		return
	}
	// 	🐉【toplink】Group Solitaire
	//  Use Solitaire to help you collect information submitted by users conveniently and quickly.

	// 接龙1
	// ├收集中
	// ├创建时间：2023-09-02 21:19:44
	// ├已收集：2条
	// └规则介绍：测试接龙1
	content := fmt.Sprintf("🐉【%s】群接龙\n使用接龙来帮你方便快捷的收集用户提交的信息。\n\n", chat.FirstName)

	for i, item := range items {
		content += fmt.Sprintf("接龙%d\n├%s\n├创建时间：%s\n├已收集：%d条\n└规则介绍：%s\n\n",
			i+1,
			solitaireStatus[item.Status],
			item.CreatedAt,
			item.MsgCollected,
			item.Description,
		)
	}
	rows := [][]model.ButtonInfo{}
	// buttons
	for i, item := range items {
		name := fmt.Sprintf("接龙%d", i+1)
		if item.Status == model.SolitaireStatusActive {
			name += "✅"
		} else {
			name += "❌"
		}
		btn1 := model.ButtonInfo{
			Text:    name,
			Data:    "solitaire_name",
			BtnType: model.BtnTypeData,
		}
		btn2 := model.ButtonInfo{
			Text:    "文件导出",
			Data:    "solitaire_export",
			BtnType: model.BtnTypeData,
		}
		btn3 := model.ButtonInfo{
			Text:    "消息导出",
			Data:    "solitaire_messages",
			BtnType: model.BtnTypeData,
		}
		btn4 := model.ButtonInfo{
			Text:    "删除",
			Data:    "solitaire_delete",
			BtnType: model.BtnTypeData,
		}
		rows = append(rows, []model.ButtonInfo{btn1, btn2, btn3, btn4})
	}
	rows = append(rows, []model.ButtonInfo{
		{
			Text:    "➕ 新建接龙",
			Data:    "solitaire_create",
			BtnType: model.BtnTypeData,
		},
	})
	rows = append(rows, []model.ButtonInfo{
		{
			Text:    "🏠 返回首页",
			Data:    "solitaire_home",
			BtnType: model.BtnTypeData,
		},
	})
	keyboard := utils.MakeKeyboard(rows)
	utils.GroupWelcomeMarkup = keyboard
	reply := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = mgr.bot.Send(reply)
	if err != nil {
		logger.Err(err).Msg("send solitaire index failed")
	}
}
