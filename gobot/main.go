package main

import (
	"flag"
	"gobot/cmd/bot"
	"gobot/internal/models"
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

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/addTaskTitle", upgradeBot.AddTaskTitle)
	upgradeBot.Bot.Handle("/addTaskDescription", upgradeBot.AddTaskDescription)
	upgradeBot.Bot.Handle("/addTaskEndDate", upgradeBot.AddTaskEndDate)
	upgradeBot.Bot.Handle("/addTask", upgradeBot.AddTask)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.StartHandler)
	// upgradeBot.Bot.Handle("/deleteTask", upgradeBot.DeleteTask)

	upgradeBot.Bot.Start()
}
