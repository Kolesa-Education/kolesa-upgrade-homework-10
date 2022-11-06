package main

import (
	"log"
	"upgrade/Model/Bot"
	"upgrade/Model/Database"
)
import "github.com/BurntSushi/toml"

type Config struct {
	BotToken   string `toml:"botToken"`
	DbAddress  string `toml:"dbAddress"`
	DbName     string `toml:"dbName"`
	DbUsername string `toml:"dbUsername"`
	DbPassword string `toml:"dbPassword"`
}

func main() {
	cfg := readConfig()
	database := Database.NewDatabase(cfg.DbAddress, cfg.DbName, cfg.DbUsername, cfg.DbPassword)
	if database.Connection.Error != nil {
		log.Fatal("Error connecting to database:", database.Connection.Error)
	}
	tgBot := Bot.Bot{
		Bot:      Bot.InitBot(cfg.BotToken),
		Database: database,
	}
	tgBot.Bot.Handle("/start", tgBot.StartHandler)
	tgBot.Bot.Handle("/tasks", tgBot.ShowTasks)
	tgBot.Bot.Handle("/newTask", tgBot.NewTask)
	tgBot.Bot.Handle("/deleteTask", tgBot.DeleteTask)
	tgBot.Bot.Start()
}

func readConfig() *Config {
	var cfg Config
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	return &cfg
}
