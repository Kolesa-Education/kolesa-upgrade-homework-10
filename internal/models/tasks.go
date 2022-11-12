package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EndDate     time.Time `json:"end_date"`
	UserId      uint      `json:"user_id"`
	User        User
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {
	result := m.Db.Create(task)

	return result.Error
}
