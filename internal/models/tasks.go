package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name       string `json:"name"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}