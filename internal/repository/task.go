package repository

import (
	"bot-tasker/internal/models"

	"gorm.io/gorm"
)

type Task interface {
	CreateTask(task models.Task) error
	GetAll(userId uint) ([]models.Task, error)
	Delete(userId, taskId uint) error
}

type TaskRepos struct {
	db *gorm.DB
}

func newTaskRepos(db *gorm.DB) *TaskRepos {
	return &TaskRepos{
		db: db,
	}
}

func (r *TaskRepos) CreateTask(task models.Task) error {
	res := r.db.Create(&task)
	return res.Error
}

func (r *TaskRepos) GetAll(userId uint) ([]models.Task, error) {
	var tasks []models.Task
	res := r.db.Find(&tasks, models.Task{UserID: userId})
	return tasks, res.Error
}

func (r *TaskRepos) Delete(userId, taskId uint) error {
	return r.db.Where("user_id=?", userId).Delete(&models.Task{}, taskId).Error
}
