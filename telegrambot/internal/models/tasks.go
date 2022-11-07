package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	TelegramUserId int64  `json:"telegram_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	EndDate        string `json:"end_date"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) Create(task Task) error {

	result := m.Db.Create(&task)
	//Create into database -- table --
	return result.Error
}

// func (m *TaskModel) GetDist(ctx id) error {
// 	//SELECT DISTINCT title description FROM table_list WHERE telegram_user_id = id

// 	result := m.Db.Pluck(id)
// 	//Create into database -- table --
// 	return result.ScanRows(id)
// }

func (m *TaskModel) GetTask(taskId int64) ([]Task, error) {
	var existTask []Task
	result := m.Db.Find(&existTask, Task{TelegramUserId: taskId})
	if result.Error != nil {
		return nil, nil
	}
	return existTask, nil
}

// DeleteTask functionx
func (m *TaskModel) DeleteUserTask(task Task) error {
	var existTask []Task
	result := m.Db.Delete(&existTask, task)

	return result.Error
}
