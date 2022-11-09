package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kolesa-upgrade-homework-10/cmd/bot"
	"kolesa-upgrade-homework-10/internal/models"
	"log"
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
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	upgradeBot := bot.UpgradeBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
	}

	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса базы данных %v", err)
	}

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Handle("/addTask", upgradeBot.AddHandler)
	upgradeBot.Bot.Handle("/tasks", upgradeBot.AllTaskHandler)
	upgradeBot.Bot.Handle("/deleteTask", upgradeBot.DeleteHandler)
	//upgradeBot.Bot.Handle("/game", upgradeBot.GameHandler)
	//upgradeBot.Bot.Handle("/try", upgradeBot.TryHandler)
	upgradeBot.Bot.Start()
}
