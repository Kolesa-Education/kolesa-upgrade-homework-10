package main

import (
	"flag"
	"kolesa-upgrade-homework-10/cmd/bot"
	"kolesa-upgrade-homework-10/internal/models"
	"log"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Env      string
	BotToken string
	Dsn      string
}

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := &Config{}
	_, err := toml.DecodeFile(*configPath, cfg)

	if err != nil {
		log.Fatalf("Ошибка декодирования файла конфигов %v", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	telegramBot := bot.TelegramBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}

	telegramBot.Bot.Handle("/start", telegramBot.Start)
	telegramBot.Bot.Handle("/help", telegramBot.Help)
	telegramBot.Bot.Handle("/addtask", telegramBot.AddTask)

	telegramBot.Bot.Handle("/tasks", telegramBot.TaskList)
	telegramBot.Bot.Handle("/delete", telegramBot.DeleteTask)

	telegramBot.Bot.Start()
}
