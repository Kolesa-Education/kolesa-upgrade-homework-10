package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EndDate     time.Time `json:"end_date"`
	ChatId      int64     `json:"chat_id"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	log.Println(&task)

	result := m.Db.Create(&task)

	return result.Error
}

func (m *TaskModel) FindOne(ChatId int64) (*Task, error) {
	existTask := Task{}

	result := m.Db.First(&existTask, User{ChatId: ChatId})

	if result.Error != nil {
		return nil, result.Error
	}

	return &existTask, nil
}

func (m *TaskModel) FindAll(chat_id int64) ([]Task, error) {
	var existTasks []Task

	result := m.Db.Find(&existTasks, Task{ChatId: chat_id})

	if result.Error != nil {
		return nil, result.Error
	}

	return existTasks, nil
}

func (m *TaskModel) Delete(taskId int64) error {
	result := m.Db.Delete(&Task{}, taskId)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
