package setting

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegramBot/model"
	"telegramBot/utils"
)

func SpamSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	btn01 := model.ButtonInfo{
		Text:    "AI屏蔽垃圾消息[强劲版]",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn11 := model.ButtonInfo{
		Text:    "反洪水攻击",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn12 := model.ButtonInfo{
		Text:    "屏蔽被封禁账号",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn21 := model.ButtonInfo{
		Text:    "屏蔽链接",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn22 := model.ButtonInfo{
		Text:    "屏蔽频道马甲发言",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn31 := model.ButtonInfo{
		Text:    "屏蔽来自频道转发",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn32 := model.ButtonInfo{
		Text:    "屏蔽来自用户转发",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn41 := model.ButtonInfo{
		Text:    "屏蔽@群组ID",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}
	btn42 := model.ButtonInfo{
		Text:    "屏蔽@用户ID",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}
	btn51 := model.ButtonInfo{
		Text:    "屏蔽以太坊地址",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}
	btn52 := model.ButtonInfo{
		Text:    "清除命令消息",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn61 := model.ButtonInfo{
		Text:    "屏蔽超长消息",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}
	btn62 := model.ButtonInfo{
		Text:    "设置超长姓名长度",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn71 := model.ButtonInfo{
		Text:    "惩罚设置",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn72 := model.ButtonInfo{
		Text:    "例外管理",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn81 := model.ButtonInfo{
		Text:    "自动删除提醒消息",
		Data:    "prohibited_ban_time",
		BtnType: model.BtnTypeData,
	}

	btn91 := model.ButtonInfo{
		Text:    "返回",
		Data:    "go_setting",
		BtnType: model.BtnTypeData,
	}
	row0 := []model.ButtonInfo{btn01}
	row1 := []model.ButtonInfo{btn11, btn12}
	row2 := []model.ButtonInfo{btn21, btn22}
	row3 := []model.ButtonInfo{btn31, btn32}
	row4 := []model.ButtonInfo{btn41, btn42}
	row5 := []model.ButtonInfo{btn51, btn52}
	row6 := []model.ButtonInfo{btn61, btn62}
	row7 := []model.ButtonInfo{btn71, btn72}
	row8 := []model.ButtonInfo{btn81}
	row9 := []model.ButtonInfo{btn91}
	rows := [][]model.ButtonInfo{row0, row1, row2, row3, row4, row5, row6, row7, row8, row9}
	keyboard := utils.MakeKeyboard(rows)
	utils.SpamSettingMenuMarkup = keyboard

	content := updateSpamMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

}

func updateSpamMsg() string {
	content := "📨 反垃圾\n\n惩罚：踢出+封禁 60 分钟\n\n自动删除提醒消息：10分钟\n\n✅AI屏蔽垃圾消息[强劲版]: \n└ 全网已拦截广告：20645283 次 查看详情\n✅ 反洪水攻击:\n└ 同一条(相似)消息一段时间内在多个群发送\n✅ 屏蔽被封禁账号:\n└ 多次发送垃圾消息被全网封禁的账号"

	return content
}
