package group

import (
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	iradix "github.com/hashicorp/go-immutable-radix/v2"
)

type CallbackParam struct {
	chatId int64
	msgId  int
	data   string // callback data
	query  string // raw query
	newMsg bool   // 是否新消息
	param  url.Values
}

type CallbackFn func(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error

var router *iradix.Tree[CallbackFn]

func NewCallbackParam(chatId int64, msgId int, data string, newmsg bool) *CallbackParam {
	cp := &CallbackParam{
		chatId: chatId,
		msgId:  msgId,
		data:   data,
		newMsg: newmsg,
	}
	if data != "" {
		ss := strings.Split(data, "?")

		if len(ss) == 2 {
			var err error
			cp.query = ss[1]
			cp.param, err = url.ParseQuery(ss[1])
			if err != nil {
				logger.Err(err).Msg("invalid data or query")
			}
		}
	}
	return cp
}

// callback entry
func CallbackHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) bool {
	cb := update.CallbackQuery
	msg := cb.Message
	chat := msg.Chat
	chatId := chat.ID

	data := cb.Data
	if data == "" {
		return false
	}

	param := &CallbackParam{
		chatId: chatId,
		data:   data,
		msgId:  msg.MessageID,
	}
	var (
		fn  CallbackFn
		ok  bool
		err error
	)
	if strings.HasSuffix(data, "$") {
		// 精确匹配
		fn, ok = router.Get([]byte(data[0 : len(data)-1]))
		if !ok {
			logger.Warn().Str("data", data).Msg("no match router")
			return false
		}
		// 无参数
		param.query = ""
		param.param = url.Values{}
	} else {
		ss := strings.Split(data, "?")
		fn, ok = router.Get([]byte(ss[0]))
		if !ok {
			return false
		}

		if len(ss) == 2 {
			param.query = ss[1]
			param.param, err = url.ParseQuery(ss[1])

			if err != nil {
				logger.Err(err).Str("data", data).Msg("invalid callback query")
				return true
			}
		} else {
			param.query = ""
			param.param = url.Values{}
		}
	}

	fn(update, bot, param)
	return true
}

func RegisterCallback(data string, fn CallbackFn) {
	// if btnType == model.BtnTypeData {
	// 	btn = tgbotapi.NewInlineKeyboardButtonData(text, data)
	// } else if btnType == model.BtnTypeUrl {
	// 	btn = tgbotapi.NewInlineKeyboardButtonURL(text, data)
	// } else if btnType == model.BtnTypeSwitch {
	// 	btn = tgbotapi.NewInlineKeyboardButtonSwitch(text, data)
	// }
	var ok bool
	if strings.HasSuffix(data, "$") {
		router, _, _ = router.Insert([]byte(data[0:len(data)-1]), fn)
	} else {
		router, _, ok = router.Insert([]byte(data), fn)
	}
	if ok {
		logger.Error().Str("route", data).Msg("insert callback route already exists")
	}
}

func InitCallbackRouters() {
	router = iradix.New[CallbackFn]()

	RegisterCallback("lucky$", luckyIndex)
	RegisterCallback("lucky_record", luckyRecords)
	RegisterCallback("lucky_cancel", luckyCancel)
	RegisterCallback("luckysetting", toggleLuckySetting)
	RegisterCallback("lucky_create_index$", LuckyCreateIndex)
	RegisterCallback("lucky_create", luckyCreate)
	RegisterCallback("lucky_create_general", luckyCreateGeneral)
	RegisterCallback("lucky_create_chatJoin", luckyCreateChatJoin)
	RegisterCallback("lucky_create_invite", luckyCreateInvite)
	RegisterCallback("lucky_create_hot", luckyCreateHot)
	RegisterCallback("lucky_create_fun", luckyCreateFun)
	RegisterCallback("lucky_create_points", luckyCreatePoints)
	RegisterCallback("lucky_create_answer", luckyCreateAnswer)
	RegisterCallback("lucky_create_keywords", luckyCreateKeywords)
	RegisterCallback("lucky_create_name", luckyCreateName)
	RegisterCallback("lucky_push", luckyCreatePush)
	RegisterCallback("lucky_publish", luckyCreatePublish)
}
