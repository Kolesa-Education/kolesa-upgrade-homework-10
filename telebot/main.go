package main

import (
	"flag"
	"log"
	"telebot/cmd/bot"
	"telebot/internal/models"

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

	teleBot := bot.TeleBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
	}

	teleBot.Bot.Handle("/start", teleBot.StartHandler)
	teleBot.Bot.Handle("/help", teleBot.HelpHandler)
	teleBot.Bot.Handle("/addTask", teleBot.AddTaskHandler)

	//teleBot.Bot.Handle("/tasks", teleBot.TasksHandler)
	//teleBot.Bot.Handle("/deleteTask {id}", teleBot.DeleteTaskHandler)

	teleBot.Bot.Start()
}
