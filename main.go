package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Kolesa-Education/kolesa-upgrade-homework-10/cmd/bot"
	"github.com/Kolesa-Education/kolesa-upgrade-homework-10/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
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

	bot := bot.Bot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TasksModel{Db: db},
	}

	bot.Bot.Handle("/start", bot.StartHandler)
	bot.Bot.Handle("/addTask", bot.AddTaskHandler)
	bot.Bot.Handle("/tasks", bot.TasksHandler)
	bot.Bot.Handle("/deleteTask", bot.DeleteTaskHandler)
	bot.Bot.Start()

	bot.Bot.Start()
}
