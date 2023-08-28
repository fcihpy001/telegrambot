package group

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"telegramBot/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 最大统计跨度
const MaxTimeRange = 86400 * 30

var (
	ErrNoMessage          = errors.New("message is nil")
	ErrNoChat             = errors.New("chat is nil")
	ErrNoReplyTo          = errors.New("replyTo message is nil")
	ErrInvalidStartTime   = errors.New("invalid start time")
	ErrInvalidTimeRange   = errors.New("invalid time range")
	ErrExceedMaxTimeRange = errors.New("exceed max time range")
	ErrInvalidTimeFormat  = errors.New("invalid time format")

	untilRe = regexp.MustCompile(`\s*(\d+)\s*(\w*)`)
	timeRe  = regexp.MustCompile(`\s*((\d{4}-)?(\d{1,2})-(\d{1,2}))?\s*((\d{1,2}):(\d{1,2})(:\d{1,2})?)?`)
)

func (mgr *GroupManager) checkAdmin(update *tgbotapi.Update) {
	msg := update.Message
	chat := msg.Chat
	replyTo := msg.ReplyToMessage
	if chat == nil || replyTo == nil {
		logger.Warn().Msg("chat is nil or replyTo is nil")
		return
	}
	mgr.CheckUserIsAdmin(chat.ID, replyTo.From.ID)
}

// get chatId, userId
// @return chatId chat id
// @return userId the replyTo message userId
// @return messageId the replyTo messageId
func getChatUserFromReplyMessage(update *tgbotapi.Update) (chatId, userId, messageId int64, err error) {
	msg := update.Message
	if msg == nil {
		err = ErrNoMessage
		return
	}
	replyTo := msg.ReplyToMessage
	chat := msg.Chat
	if chat == nil {
		err = ErrNoChat
		return
	}
	if replyTo == nil {
		err = ErrNoReplyTo
		return
	}
	chatId = chat.ID
	userId = replyTo.From.ID
	messageId = int64(replyTo.MessageID)
	return
}

// 用于 Ban Mute 指令
// 4分钟 4
// 4分钟 4m 4min 4minute 4minutes
// 4小时 4h 4hour 4hours
// 4天 4d 4day 4days
// 4周 4w 4week 4weeks
// 4月 4month 4months
// 4年 4y 4year
func parseUntilDate(s string) int64 {
	return time.Now().Unix() + convertToSeconds(s)
}

func convertToSeconds(s string) int64 {
	mats := untilRe.FindStringSubmatch(s)
	if len(mats) != 3 {
		logger.Warn().Msgf("invalid until pattern: %s\n", s)
		return 0
	}
	num, _ := strconv.Atoi(mats[1])

	unit := strings.TrimSpace(strings.ToLower(mats[2]))
	coefficient := 1
	switch unit {
	case "":
		fallthrough
	case "m":
		fallthrough
	case "min":
		fallthrough
	case "minute":
		fallthrough
	case "minutes":
		coefficient = 60

	case "h":
		fallthrough
	case "hour":
		fallthrough
	case "hours":
		coefficient = 3600

	case "d":
		fallthrough
	case "day":
		fallthrough
	case "days":
		coefficient = 86400

	case "w":
		fallthrough
	case "week":
		fallthrough
	case "weeks":
		coefficient = 86400 * 7

	case "month":
		fallthrough
	case "months":
		coefficient = 86400 * 30

	case "y":
		fallthrough
	case "year":
		fallthrough
	case "years":
		coefficient = 86400 * 365
	}
	return int64(num * coefficient)
}

//lint:ignore U1000 ignore unused
func prettyJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "  ")
	return string(buf)
}

// {{yyyy-}mm-dd} hh:MM{:SS} | {{yyyy-}mm-dd} hh:MM{:SS}
// @return startTs start timestamp in second
// @return endTs end timestamp in second
// @return error error
func parseTimeRange(s string) (startTs int64, endTs int64, err error) {
	s = strings.Replace(s, "/stats", "", -1)
	s = strings.TrimSpace(s)
	if s == "" {
		// 今天
		now := time.Now().Unix()
		return now - now%86400, now, nil
	}

	ss := strings.Split(s, "|")
	if len(ss) != 2 {
		return
	}
	if startTs, err = parseTime(ss[0], true); err != nil {
		return
	}
	if endTs, err = parseTime(ss[1], false); err != nil {
		return
	}
	now := time.Now().Unix()
	if startTs >= now {
		err = ErrInvalidStartTime
		return
	}
	if endTs > now {
		endTs = now
	}
	if endTs-startTs > MaxTimeRange {
		err = ErrExceedMaxTimeRange
		return
	}
	return
}

// start: the first time of time range
func parseTime(s string, start bool) (int64, error) {
	mats := timeRe.FindStringSubmatch(s)
	// println("mats:", mats)
	if len(mats) != 9 {
		logger.Warn().Msgf("invalid time range pattern: %s\n", s)
		return 0, ErrInvalidTimeFormat
	}
	now := time.Now()
	// mats[0] whole pattern
	// mats[1] yyyy-mm-dd
	// mats[2] yyyy-
	// mats[3] mm
	// mats[4] dd
	// mats[5] hh:MM:ss
	// mats[6] hh
	// mats[7] MM
	// mats[8] :ss
	var (
		yyyy  int
		month int
		day   int
		hour  int
		mm    int
		ss    int
	)
	if mats[1] == "" {
		yyyy = now.Year()
		month = int(now.Month())
		day = now.Day()
	} else {
		if mats[2] == "" {
			yyyy = now.Year()
		} else {
			yyyy, _ = strconv.Atoi(mats[2][0 : len(mats[2])-1])
		}
		month, _ = strconv.Atoi(mats[3])
		day, _ = strconv.Atoi(mats[4])
	}

	if mats[5] == "" {
		if start {
			hour, mm, ss = 0, 0, 0
		} else {
			hour, mm, ss = 23, 59, 59
		}
	} else {
		hour, _ = strconv.Atoi(mats[6])
		mm, _ = strconv.Atoi(mats[7])
		ss, _ = strconv.Atoi(mats[8])
	}

	// todo time.Local or UTC? or something else?
	tm := time.Date(yyyy, time.Month(month), day, hour, mm, ss, 0, time.Local)
	return tm.Unix(), nil
}

func (mgr *GroupManager) sendMessage(c tgbotapi.Chattable, fmt string, args ...interface{}) {
	utils.SendMessage(mgr.bot, c, fmt, args...)
}

func (mgr *GroupManager) sendText(chatId int64, text string) {
	utils.SendText(chatId, text, mgr.bot)
}
