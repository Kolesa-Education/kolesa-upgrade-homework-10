package models

import (
	"gorm.io/gorm"
	"log"
)

type Task struct {
	gorm.Model
	Title   string `json:"title"`
	Descr   string `json:"descr"`
	EndDate string `json:"endDate"`
	Userid  int64  `json:"userId"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) CreateTask(task Task) error {
	result := m.Db.Create(&task)
	return result.Error
}

func (m *TaskModel) ShowTaskDb(userId int64) (error, Task) {
	tasks := Task{}
	err := m.Db.Preload("Tasks").Find(&tasks, userId)
	if err.Error != nil {
		log.Print("Fatal")
	}
	log.Fatal(err)
	return nil, tasks
}
