package main

import (
	"log"

	"bot-tasker/internal/app"
	"bot-tasker/internal/config"
)

func main() {
	log.Println("init config")
	cfg, err := config.GetConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	app := app.NewApp(cfg)
	app.Run()
}
