package service

import (
	"bot-tasker/internal/models"
	"bot-tasker/internal/repository"
)

type Task interface {
	CreateTask(task models.Task) error
	Delete(user *models.User, taskId uint) error
	GetAll(userId uint) ([]models.Task, error)
}

type TaksService struct {
	repo repository.Task
}

func newTaskService(r repository.Task) *TaksService {
	return &TaksService{
		repo: r,
	}
}

func (s *TaksService) CreateTask(task models.Task) error {
	return s.repo.CreateTask(task)
}

func (s *TaksService) GetAll(userId uint) ([]models.Task, error) {
	return s.repo.GetAll(userId)
}

func (s *TaksService) Delete(user *models.User, taskId uint) error {
	return s.repo.Delete(user.ID, taskId)
}
