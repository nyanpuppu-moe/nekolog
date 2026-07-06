package service

import (
	"nekolog/internal/model"
	"nekolog/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) Get(name string) (*model.User, error) {
	user, err := s.userRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindUserByID(id model.UserID) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
