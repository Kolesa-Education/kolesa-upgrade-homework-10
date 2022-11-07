package repository

import (
	"upgrade/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Db: db,
	}
}

func (r *Repository) CreateUser(user models.User) error {
	result := r.Db.Create(&user)

	return result.Error
}

func (r *Repository) FindUser(telegramId int64) (*models.User, error) {
	var user *models.User

	result := r.Db.First(&user, models.User{TelegramId: telegramId})

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *Repository) CreateTask(task models.Task) error {
	result := r.Db.Create(&task)

	return result.Error
}

func (r *Repository) FindTask(id int) (*models.Task, error) {
	var task *models.Task

	result := r.Db.First(&task, models.Task{Id: id})

	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (r *Repository) AllTasks(user *models.User) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.Db.Model(&user).Association("Tasks").Find(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) DeleteTask(id int) {
	r.Db.Delete(&models.Task{}, id)
}
