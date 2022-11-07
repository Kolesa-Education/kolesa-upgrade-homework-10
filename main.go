package main

import (
	"flag"
	"log"

	"upgrade/cmd/bot"
	"upgrade/internal/repository"

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

	taskBot := bot.TaskBot{
		Bot:  bot.InitBot(cfg.BotToken),
		Repo: repository.NewRepository(db),
	}

	taskBot.InitCommands()
	// taskBot.Bot.Handle("/start", taskBot.Start)
	// taskBot.Bot.Handle("/addTask", taskBot.AddTask)
	taskBot.Bot.Start()
}
