package Bot

import (
	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"log"
	"time"
	"upgrade/Model/Database"
)

type Bot struct {
	Database *gorm.DB
	Bot      *telebot.Bot
}

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return b
}

func (bot *Bot) StartHandler(ctx telebot.Context) error {
	Database.NewUser(bot.Database, ctx.Sender().ID, ctx.Sender().FirstName, ctx.Sender().LastName, ctx.Chat().ID)
	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}
