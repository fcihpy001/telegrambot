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
	param  url.Values
}

type CallbackFn func(update *tgbotapi.Update, bot *tgbotapi.BotAPI, param *CallbackParam) error

var router *iradix.Tree[CallbackFn]

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
	RegisterCallback("lucky_create_index$", luckyCreateIndex)
	RegisterCallback("lucky_create", luckyCreate)
	RegisterCallback("lucky_create_general", luckyCreateGeneral)
	RegisterCallback("lucky_create_keywords", luckyCreateKeywords)
	RegisterCallback("lucky_push", luckyCreatePush)
	RegisterCallback("lucky_publish", luckyCreatePublish)
}
