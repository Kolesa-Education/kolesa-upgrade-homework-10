package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
	UserId      uint   `json:"user_id"`
}
