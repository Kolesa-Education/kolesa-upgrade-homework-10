package main

import (
	"flag"
	"log"
	"telegrambot/cmd/bot"
	"telegrambot/internal/models"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Env      string
	BotToken string
	Dsn      string
	Dsn2     string
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
	db2, err := gorm.Open(sqlite.Open(cfg.Dsn2), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db2},
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/addTask", upgradeBot.AddTaskHandler)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.TasksHandler)
	// upgradeBot.Bot.Handle("/deleteTask{id}", upgradeBot.DeleteTaskHandler)

	// upgradeBot.Bot.Handle("/game", upgradeBot.GameHandler)
	upgradeBot.Bot.Handle("/add", upgradeBot.AddHandler)
	upgradeBot.Bot.Handle("/delete", upgradeBot.DeleteTaskHandler)

	upgradeBot.Bot.Start()
}
