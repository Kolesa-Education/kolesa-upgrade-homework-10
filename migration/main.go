package main

import (
	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kolesa-upgrade-hw-10/internal/models"
	"log"
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

	//db.Table("tasks", models.Task{}).Col
	db.AutoMigrate(&models.User{}, &models.Task{})
}
