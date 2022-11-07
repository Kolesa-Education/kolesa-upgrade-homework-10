package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string    `json:"title"`
	UserId      int64     `json:"user_id"`
	Description string    `json:"description"`
	EndDate     time.Time `json:"end_date"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) FindTasks(telegramId int64) ([]Task, error) {
	existTasks := []Task{}

	result := m.Db.Find(&existTasks, Task{UserId: telegramId})

	if result.Error != nil {
		return nil, result.Error
	}

	return existTasks, nil
}
