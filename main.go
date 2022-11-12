package main

import (
	"flag"
	"log"
	"tasking/cmd/bot"
	"tasking/internal/models"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Env      string
	BotToken string
}

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := &Config{}
	_, err := toml.DecodeFile(*configPath, cfg)

	if err != nil {
		log.Fatalf("Ошибка декодирования файла конфигов %v", err)
	}

	dsn := "root:Miracle.1208@tcp(127.0.0.1:3307)/tasking?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/addTask", upgradeBot.AddTaskHandler)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.TasksHandler)
	upgradeBot.Bot.Handle("/delete", upgradeBot.DeleteTaskHandler)

	upgradeBot.Bot.Start()
}
