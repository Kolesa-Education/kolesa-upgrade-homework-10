package main

import (
	"flag"
	"log"
	"taskbot/cmd/bot"
	"taskbot/internal/models"

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
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
		Tasks: &models.TaskModel{Db: db},
	}

	taskBot.Bot.Handle("/start", taskBot.StartHandler)
	taskBot.Bot.Handle("/game", taskBot.GameHandler)
	taskBot.Bot.Handle("/try", taskBot.TryHandler)
	taskBot.Bot.Handle("/taskrule", taskBot.TaskRuleHandler)
	taskBot.Bot.Handle("/addtask", taskBot.TaskHandler)
	taskBot.Bot.Handle("/tasks", taskBot.AllTasksHandler)

	taskBot.Bot.Start()
}
