package app

import (
	"log"
	"time"

	"bot-tasker/internal/config"
	"bot-tasker/internal/delivery"
	"bot-tasker/internal/models"
	"bot-tasker/internal/repository"
	"bot-tasker/internal/service"

	"gopkg.in/telebot.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	bot    *telebot.Bot
}

func NewApp(cfg *config.Config) *App {
	log.Println("init telebot")
	pref := telebot.Settings{
		Token:  cfg.Bot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		config: cfg,
		bot:    b,
	}
}

func (a *App) Run() {
	log.Println("init gorm db")
	db, err := gorm.Open(sqlite.Open(a.config.Gorm.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(models.User{}, models.Task{}); err != nil {
		log.Fatal(err)
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := delivery.NewHandler(service)

	a.bot.Handle("/start", handler.StartHandler)
	a.bot.Handle("/addTask", handler.AddTaskHandler)
	a.bot.Handle("/tasks", handler.AllTasksHandler)
	a.bot.Handle("/deleteTask", handler.DeleteTaskHandler)
	a.bot.Handle(telebot.OnText, handler.NewMsgHandler)

	log.Println("start bot")
	a.bot.Start()
}
