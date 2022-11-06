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

/*func (m *TaskModel) GetAll(taskId int) (*Task, error) {

	tasks := Task{}

	err := m.Db.Preload("Tasks").Find(&tasks, taskId).Error

	return &tasks, err
}

func (m *TaskModel) DeleteTask(taskId string, userId int64) bool {

	m.Db.Where("user_id = ?", userId).Where("id = ?", taskId).Delete(&Task{})

	return true
}*/
