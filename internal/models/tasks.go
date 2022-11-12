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
	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) DeleteTask(taskId int64) error {
	error := m.Db.First(&Task{}, taskId).Delete(&Task{}).Error

	return error
}

func (m *TaskModel) GetAll(userId uint) ([]Task, error) {
	tasks := make([]Task, 0)
	result := m.Db.Find(&tasks, Task{UserId: userId})

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}
