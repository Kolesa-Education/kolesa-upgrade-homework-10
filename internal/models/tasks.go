package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	ID    int `json:"id"`
	Title       string `json:"title"`
	Description string  `json:"description"`
	UserId      uint  `json:"user_id"`
	EndDate  time.Time `json:"end_date"`
	User User `json:"user"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) GetAllUserTasks(userId int) (*[]Task, error) {
	var tasks []Task

	result := m.Db.Where("user_id = ?", userId).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tasks, nil
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) DeleteUserTask(task Task, userId int) error {
	result := m.Db.Where("user_id = ?", userId).Delete(&task)

	return result.Error
}

func (m *TaskModel) FindOne(taksId int) (*Task, error) {
	existTask := Task{}

	result := m.Db.First(&existTask, taksId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &existTask, nil
}

func (m *TaskModel) FindOneForUser(taksId int, userId int) (*Task, error) {
	existTask := Task{}

	result := m.Db.Where("user_id = ?", userId).First(&existTask, taksId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &existTask, nil
}