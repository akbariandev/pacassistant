package telegram

import (
	telegram "gopkg.in/telebot.v3"
	"time"
)

func NewTelegramBot(token string) *telegram.Bot {
	pref := telegram.Settings{
		Token:  token,
		Poller: &telegram.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telegram.NewBot(pref)
	if err != nil {
		panic(err)
	}

	return bot
}
