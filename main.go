package main

import (
	"kolesa-upgrade-hw-10/cmd/bot"
	"kolesa-upgrade-hw-10/internal/models"
	"log"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Env      string
	BotToken string
	Dsn      string
}

func main() {
	cfg := &Config{}
	_, err := toml.DecodeFile("config/local.toml", cfg)

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
	upgradeBot.Bot.Handle("/addtask", upgradeBot.CreateTaskHandler)
	upgradeBot.Bot.Handle("/showtasks", upgradeBot.ShowTasksHandler)
	upgradeBot.Bot.Handle("/deletetask", upgradeBot.DeleteTaskHandler)
	upgradeBot.Bot.Handle("/help", upgradeBot.HelpHandler)

	upgradeBot.Bot.Start()
}
