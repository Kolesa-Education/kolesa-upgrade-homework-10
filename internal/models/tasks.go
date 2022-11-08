package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserId      int64     `json:"user_id"`
	EndDate     time.Time `json:"end_date"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) DeleteTask(task_id int64, user_id int64) error {
	db := m.Db.Where("user_id = ?", user_id).Where("id = ?", task_id).Delete(&Task{})

	return db.Error
}

func (m *TaskModel) FindAll(userId int64) ([]Task, error) {
	var tasks []Task

	result := m.Db.Find(&tasks, Task{UserId: userId})

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (m *TaskModel) Create(task Task) error {
	result := m.Db.Create(&task)
	return result.Error
}
