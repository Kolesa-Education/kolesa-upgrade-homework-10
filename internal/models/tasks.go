package models

import (
	"fmt"
	"gorm.io/gorm"
)

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

func (m *TaskModel) FindAll(userId uint) ([]Task, error) {
	var existTasks []Task

	result := m.Db.Find(&existTasks, Task{UserID: userId})
	fmt.Println(existTasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return existTasks, nil
}

func (m *TaskModel) DeleteOne(taskID int) error {
	existTask := Task{}

	result := m.Db.Delete(&existTask, taskID)

	return result.Error
}
