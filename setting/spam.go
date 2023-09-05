package setting

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"telegramBot/model"
	"telegramBot/services"
	"telegramBot/utils"
)

var (
	spamsSetting = model.SpamSetting{}
)

func spamSettingMenu(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	err = services.GetModelData(utils.GroupInfo.GroupId, &spamsSetting)
	fmt.Println("spamsSetting-query", spamsSetting)
	spamsSetting.ChatId = utils.GroupInfo.GroupId

	var buttons [][]model.ButtonInfo
	utils.Json2Button2("spam.json", &buttons)
	fmt.Println(&buttons)
	var rows [][]model.ButtonInfo
	for i := 0; i < len(buttons); i++ {
		btnArr := buttons[i]
		var row []model.ButtonInfo
		for j := 0; j < len(btnArr); j++ {
			btn := btnArr[j]
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	keyboard := utils.MakeKeyboard(rows)
	utils.SpamSettingMenuMarkup = keyboard

	content := updateSpamMsg()
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, content, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

}

func SpamSettingHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := update.CallbackQuery.Data
	query := strings.Split(data, ":")
	cmd := query[0]
	params := ""
	if len(query) > 1 {
		params = query[1]
	}

	if cmd == "spam_setting" {
		spamSettingMenu(update, bot)

	} else if cmd == "spam_setting_type" {
		typeStatusHandler(update, bot, params)

	} else if cmd == "spam_setting_msg_length" {
		msgLengthHandler(update, bot)

	} else if cmd == "spam_setting_name_length" {
		nameLengthHandler(update, bot)
	}
}

func typeStatusHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI, params string) {
	if len(params) == 0 {
		return
	}

	switch params {
	case "ai":
		spamsSetting.EnableAi = !spamsSetting.EnableAi
		if spamsSetting.EnableAi {
			utils.SpamSettingMenuMarkup.InlineKeyboard[0][0].Text = "✅AI屏蔽垃圾消息[强劲版]"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[0][0].Text = "AI屏蔽垃圾消息[强劲版]"
		}

	case "ddos":
		spamsSetting.DDos = !spamsSetting.DDos
		if spamsSetting.DDos {
			utils.SpamSettingMenuMarkup.InlineKeyboard[1][0].Text = "✅反洪水攻击"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[1][0].Text = "反洪水攻击"
		}
	case "blackUser":
		spamsSetting.BlackUser = !spamsSetting.BlackUser
		if spamsSetting.BlackUser {
			utils.SpamSettingMenuMarkup.InlineKeyboard[1][1].Text = "✅屏蔽被封禁账号"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[1][1].Text = "屏蔽被封禁账号"
		}
	case "link":
		spamsSetting.Link = !spamsSetting.Link
		if spamsSetting.Link {
			utils.SpamSettingMenuMarkup.InlineKeyboard[2][0].Text = "✅屏蔽链接"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[2][0].Text = "屏蔽链接"
		}
	case "channelCopy":
		spamsSetting.ChannelCopy = !spamsSetting.ChannelCopy
		if spamsSetting.ChannelCopy {
			utils.SpamSettingMenuMarkup.InlineKeyboard[2][1].Text = "✅屏蔽频道马甲发言"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[2][1].Text = "屏蔽频道马甲发言"
		}
	case "channelForward":
		spamsSetting.ChannelForward = !spamsSetting.ChannelForward
		if spamsSetting.ChannelForward {
			utils.SpamSettingMenuMarkup.InlineKeyboard[3][0].Text = "✅屏蔽来自频道转发"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[3][0].Text = "屏蔽来自频道转发"
		}
	case "userForward":
		spamsSetting.UserForward = !spamsSetting.UserForward
		if spamsSetting.UserForward {
			utils.SpamSettingMenuMarkup.InlineKeyboard[3][1].Text = "✅屏蔽来自用户转发"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[3][1].Text = "屏蔽来自用户转发"
		}
	case "atGroup":
		spamsSetting.AtGroup = !spamsSetting.AtGroup
		if spamsSetting.AtGroup {
			utils.SpamSettingMenuMarkup.InlineKeyboard[4][0].Text = "✅屏蔽@群组ID"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[4][0].Text = "屏蔽@群组ID"
		}
	case "atUser":
		spamsSetting.AtUser = !spamsSetting.AtUser
		if spamsSetting.AtUser {
			utils.SpamSettingMenuMarkup.InlineKeyboard[4][1].Text = "✅屏蔽@用户ID"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[4][1].Text = "屏蔽@用户ID"
		}
	case "ethAddress":
		spamsSetting.EthAddr = !spamsSetting.EthAddr
		if spamsSetting.EthAddr {

			utils.SpamSettingMenuMarkup.InlineKeyboard[5][0].Text = "✅屏蔽以太坊地址"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[5][0].Text = "屏蔽以太坊地址"
		}
	case "command":
		spamsSetting.Command = !spamsSetting.Command
		if spamsSetting.Command {
			utils.SpamSettingMenuMarkup.InlineKeyboard[5][1].Text = "✅清除命令消息"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[5][1].Text = "清除命令消息"
		}
	case "longMsg":
		spamsSetting.LongMsg = !spamsSetting.LongMsg
		if spamsSetting.LongMsg {
			utils.SpamSettingMenuMarkup.InlineKeyboard[6][0].Text = "✅屏蔽超长消息"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[6][0].Text = "屏蔽超长消息"
		}
	case "longName":
		spamsSetting.LongName = !spamsSetting.LongName
		if spamsSetting.LongName {
			utils.SpamSettingMenuMarkup.InlineKeyboard[7][0].Text = "✅屏蔽超长名字"
		} else {
			utils.SpamSettingMenuMarkup.InlineKeyboard[7][0].Text = "屏蔽超长名字"
		}

	}

	updateSpamMsg()
	editText := tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		utils.SpamSettingMenuMarkup,
	)
	bot.Send(editText)
}

func msgLengthHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("📨 反垃圾\n\n检测到消息内容长度大于设定数时，将会判定为超长消息，并作出相应处罚\n\n当前设置最大长度：%d\n\n👉 输入允许的消息最大长度（例如：100）：", spamsSetting.MsgLength)
	sendReplyMsg(update.CallbackQuery.Message.Chat.ID, content, bot)
}

func nameLengthHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	content := fmt.Sprintf("📨 反垃圾\n\n检测到姓名长度大于设定数时，将会判定为超长姓名，并作出相应处罚\n\n当前设置最大长度：%d\n\n👉 输入允许的姓名最大长度（例如：15）：", spamsSetting.MsgLength)
	sendReplyMsg(update.CallbackQuery.Message.Chat.ID, content, bot)
}

func SpamNameLengthReply(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	length, _ := strconv.Atoi(update.Message.Text)
	spamsSetting.NameLength = length

	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "spam_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)
	updateSpamMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅设置成功，点击按钮返回.")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func SpamMsgLengthReply(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	length, _ := strconv.Atoi(update.Message.Text)
	spamsSetting.MsgLength = length
	btn1 := model.ButtonInfo{
		Text:    "返回",
		Data:    "spam_setting",
		BtnType: model.BtnTypeData,
	}
	row1 := []model.ButtonInfo{btn1}
	rows := [][]model.ButtonInfo{row1}
	keyboard := utils.MakeKeyboard(rows)
	updateSpamMsg()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅设置成功，点击按钮返回.")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func sendReplyMsg(chatId int64, content string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, content)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("返回"),
		))

	msg.ReplyMarkup = keyboard
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
	}
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func updateSpamMsg() string {
	content := "📨 反垃圾\n\n惩罚：踢出+封禁 60 分钟\n\n自动删除提醒消息：10分钟\n\n✅AI屏蔽垃圾消息[强劲版]: \n└ 全网已拦截广告：20645283 次 查看详情\n✅ 反洪水攻击:\n└ 同一条(相似)消息一段时间内在多个群发送\n✅ 屏蔽被封禁账号:\n└ 多次发送垃圾消息被全网封禁的账号"
	//spamsSetting.Punishment.PunishType = model.PunishTypeBan
	//spamsSetting.Punishment.BanTime = 60
	//spamsSetting.Punishment.WarningCount = 3

	services.SaveModel(&spamsSetting, spamsSetting.ChatId)
	return content
}
