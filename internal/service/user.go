package service

import (
	"errors"

	"bot-tasker/internal/models"
	"bot-tasker/internal/repository"

	"gorm.io/gorm"
)

type User interface {
	AddNewUser(user models.User) error
	Get(userID int64) (*models.User, error)
}

type UserService struct {
	repo repository.User
}

func newUserService(r repository.User) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) Get(userID int64) (*models.User, error) {
	return s.repo.GetUserByChatId(userID)
}

func (s *UserService) AddNewUser(user models.User) error {
	tempUser, err := s.repo.GetUserByChatId(user.ChatId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if tempUser == nil {
		if err := s.repo.AddUser(user); err != nil {
			return err
		}
	}
	return nil
}
