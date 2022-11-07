package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"upgrade/cmd/bot"
	"upgrade/internal/models"
)

type Config struct {
	Env      string
	BotToken string
	Dsn      string
}

func main() {
	configPath := "./config/local.toml"
	flag.Parse()

	cfg := &Config{}
	_, err := toml.DecodeFile(configPath, cfg)

	if err != nil {
		log.Fatalf("Ошибка декодирования файла конфигов %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/game", upgradeBot.GameHandler)
	upgradeBot.Bot.Handle("/try", upgradeBot.TryHandler)
	upgradeBot.Bot.Handle("/addtask", upgradeBot.AddTaskHandler)
	upgradeBot.Bot.Handle("/newtask", upgradeBot.NewTaskHandler)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.TasksHandler)
	upgradeBot.Bot.Handle("/deleteTask", upgradeBot.DeleteTaskHandler)

	upgradeBot.Bot.Start()
}
