package bot

import (
	"gopkg.in/telebot.v3"
	"log"
	"time"
)

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatalf("Ошибка инициализации бота %v", err)
	}

	return bot
}
