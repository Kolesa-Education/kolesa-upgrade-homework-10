package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EndDate     time.Time `json:"end_date"`
	TelegramID  int64     `json:"telegram_id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) AddTask(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}
