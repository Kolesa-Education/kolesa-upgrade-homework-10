package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	End_date    string `json:"end_date"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {
	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) FindSame(title, description string) (*Task, error) {
	existTask := Task{}

	result := m.Db.First(&existTask, Task{Title: title, Description: description})

	if result.Error != nil {
		return nil, result.Error
	}

	return &existTask, nil
}
