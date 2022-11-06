package main

import (
	"ZakirAvrora/homework-10/cmd/bot"
	"ZakirAvrora/homework-10/internal/models"
	"flag"
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
	configPath := flag.String("config", "./config/local.toml", "Path to config file")
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

	log.Println("Connection Established")

	db.AutoMigrate(models.User{}, models.Task{})

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/help", upgradeBot.GameHandler)
	upgradeBot.Bot.Handle("/addTask", upgradeBot.AddTaskHandler)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.TasksHandler)
	upgradeBot.Bot.Handle("/deleteTask", upgradeBot.DeleteTaskHandler)

	upgradeBot.Bot.Start()
}
