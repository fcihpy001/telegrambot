package group

import (
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserPhotos(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	mgr := GroupManager{bot}

	fixtures := []struct {
		uid      int64
		hasPhoto bool
		err      error
	}{
		{6450102772, true, nil},
		{5394405541, false, nil},
		{6616020782, true, nil},
	}

	for _, item := range fixtures {
		hasPhoto, err := mgr.HasUserProfilePhotos(item.uid)
		assert.Equal(t, item.err, err, "user: %d", item.uid)
		assert.Equal(t, item.hasPhoto, hasPhoto, "user: %d", item.uid)
	}
}

func TestDeleteMessage(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	mgr := GroupManager{bot}

	err = mgr.deleteMessage(-1001631341126, 28)
	assert.Nil(t, err)
}
