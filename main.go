package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"upgrade/cmd/bot"
	"upgrade/internal/models"
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
		log.Fatalf("Ошибка декодирования файла конфига %v", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	err = db.AutoMigrate(models.User{}, models.Task{})
	if err != nil {
		log.Fatalf("Ошибка миграции %v", err)
	}

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}
	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	//upgradeBot.Bot.Handle("/addTask", upgradeBot.AddTaskHandler)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.TasksHandler)
	upgradeBot.Bot.Start()
}
