package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
	UserID      uint   `json:"user_id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) FindAll(userId int64) (*Task, error) {
	tasks := Task{}

	result := m.Db.Preload("Tasks").Find(&tasks, userId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tasks, nil
}
