package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
	UserID      int64  `json:"user_id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) GetAll(userId int64) ([]Task, error) {

	tasks := []Task{}
	result := m.Db.Where("user_id = ?", userId).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (m *TaskModel) DeleteTask(taskId int) error {

	task := Task{}

	result := m.Db.Where("id = ?", taskId).Delete(&task, taskId)

	return result.Error
}
