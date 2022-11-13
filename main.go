package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
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
	configPath := flag.String("config", "config/local.toml", "Path to config file")
	flag.Parse()
	cfg := &Config{}
	_, err := toml.DecodeFile(*configPath, cfg)
	if err != nil {
		log.Fatalf("Ошибка чтения файла конфигурации %v", err)
	}
	db, err := gorm.Open(sqlite.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных %v", err)
	}
	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/help", upgradeBot.HelpHandler)
	upgradeBot.Bot.Handle("/show", upgradeBot.ShowTaskHandler)
	upgradeBot.Bot.Handle("/new", upgradeBot.NewTaskHandler)
	upgradeBot.Bot.Handle("/delete", upgradeBot.DeleteTaskHandler)
	upgradeBot.Bot.Start()
}
