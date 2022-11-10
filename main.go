package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"telegrambot/cmd/bot"
	"telegrambot/internal/models"
)

type Config struct {
	Env      string
	BotToken string
	Dsn      string
}

func main() {
	configPath := flag.String("config", "", "Path the config file")
	flag.Parse()
	conf := &Config{}
	_, err := toml.DecodeFile(*configPath, conf)
	if err != nil {
		log.Fatalf("Ошибка декодирования файла конфигов %v", err)
	}

	db, err := gorm.Open(sqlite.Open(conf.Dsn), &gorm.Config{})

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(conf.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/addTask", upgradeBot.TaskHandler)
	upgradeBot.Bot.Handle("/task", upgradeBot.AddTaskHandler)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.AllTaskHandler)
	upgradeBot.Bot.Handle("/deleteTask", upgradeBot.DeleteTaskHandler)

	upgradeBot.Bot.Start()
}
