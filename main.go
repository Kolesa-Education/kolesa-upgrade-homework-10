package main

import (
	"log"
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
	//Model.NewUser(database, "TGID5", "FN5", "LN5", 111)
	//Model.NewTask(database, 4, "Title6", "Desc6", time.Now())
	//Database.GetUser(database, 5)
	Database.GetTask(database, 5)
	/*tgBot := Bot.Bot{
		Bot:      Bot.InitBot(cfg.BotToken),
		Database: Database.NewDatabase(cfg.DbAddress, cfg.DbName, cfg.DbUsername, cfg.DbPassword),
	}
	tgBot.Bot.Handle("/start", tgBot.StartHandler)
	tgBot.Bot.Start()*/
}

func readConfig() *Config {
	var cfg Config
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	return &cfg
}
