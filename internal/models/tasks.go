package models

import (
    "gorm.io/gorm"
    "time"
    "log"
)

type Task struct {
    gorm.Model
    Name       string `json:"name"`
    EndDate time.Time `json:"endTime"`
    UserID uint `json:"userID"`

    
}

type TaskModel struct {
    Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {
    
    result := m.Db.Create(&task)
    log.Printf(task.Name)

    return result.Error
}

func (m *TaskModel) FindOne(userID uint) (*Task, error) {
    existTask := Task{}
    

    result := m.Db.First(&existTask, Task{UserID: userID})

    if result.Error != nil {
        return nil, result.Error
    }

    return &existTask, nil
}

func (m *TaskModel) FindByUserId(telegramId int64) (*[]Task, error) {
	var tasks []Task
	result := m.Db.Find(&tasks, "user_id = ?", telegramId)
	return &tasks, result.Error
}

