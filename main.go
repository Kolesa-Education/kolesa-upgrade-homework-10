package main

import (
    "flag"
    "log"
    "upgrade/cmd/bot"
    "upgrade/internal/models"

    "github.com/BurntSushi/toml"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "gopkg.in/telebot.v3"
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

    state := models.State{
        Storage: make(map[string]string),
    }

    addTaskState := models.AddTaskState{
        CurrentState: models.None,
        State: state,
    }
    
    upgradeBot := bot.UpgradeBot{
        Bot:   bot.InitBot(cfg.BotToken),
        Users: &models.UserModel{Db: db},
        Tasks: &models.TaskModel{Db: db},
        AddTaskState: addTaskState,
        
    }

    upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
    upgradeBot.Bot.Handle("/addTask", upgradeBot.AddTaskHandler)
    upgradeBot.Bot.Handle("/getTasks", upgradeBot.GetTasksHandler)
    upgradeBot.Bot.Handle("/deleteTask", upgradeBot.DeleteTaskHandler)
    upgradeBot.Bot.Handle(telebot.OnText, upgradeBot.GeneralHandler)

    upgradeBot.Bot.Start()
}
