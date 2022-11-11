package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Id          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	EndDate     datatypes.Date `json:"end_date"`
	UserId      int            `gorm:"column:user_id"`
	// User        *models.UserModel
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) GetTasks(user_id int) ([]Task, error) {
	tasks := []Task{}

	result := m.Db.Find(&tasks, Task{UserId: user_id})

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (m *TaskModel) Delete(task_id int) error {
	tasks := []Task{}

	result := m.Db.Delete(&tasks, Task{Id: task_id})

	return result.Error
}
