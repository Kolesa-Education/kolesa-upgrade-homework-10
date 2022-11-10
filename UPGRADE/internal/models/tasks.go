package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Task        string    `json:"task"`
	Description string    `json:"description"`
	End_date    time.Time `json:"end_date"`
	TelegramId  int64     `json:"telegram_id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}
func (m TaskModel) AllTask(telegram_id int64) ([]Task, error) {
	var tasks []Task
	db := m.Db.Find(&tasks, "telegram_id = ?", telegram_id)
	return tasks, db.Error
}

func (m *TaskModel) DeleteTask(task_id int, telegram_id int64) error {
	db := m.Db.Where("telegram_id = ?", telegram_id).Where("id = ?", task_id).Delete(&Task{})
	return db.Error
}
