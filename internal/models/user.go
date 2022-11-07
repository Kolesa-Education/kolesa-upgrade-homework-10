package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string
	TelegramId int64
	FirstName  string
	LastName   string
	ChatId     int64
	Tasks      []Task
}
