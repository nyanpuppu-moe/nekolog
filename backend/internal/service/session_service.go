package service

import (
	"errors"

	"nekolog/internal/dto"
	"nekolog/internal/model"
	"nekolog/internal/repository"
	"nekolog/internal/utils"
)

type SessionService struct {
	userRepo *repository.UserRepository
}

func NewSessionService(repo *repository.UserRepository) *SessionService {
	return &SessionService{userRepo: repo}
}

func (s *SessionService) Register(req dto.SessionRegisterRequest) error {
	reservedNames := map[string]bool{
		"login":    true,
		"logout":   true,
		"register": true,
		"api":      true,
		"admin":    true,
		"static":   true,
		"public":   true,
	}

	if reservedNames[req.Name] {
		return errors.New("Can not use this name")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Password:    hashedPassword,
	}

	return s.userRepo.Create(user)
}

func (s *SessionService) Login(req dto.SessionLoginRequest) (*model.User, error) {
	user, err := s.userRepo.FindByName(req.Name)

	if err != nil {
		return nil, errors.New("Incorrect email or password")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if utils.ComparePassword(user.Password, hashedPassword) {
		return nil, errors.New("Incorrect email or password")
	}

	return user, nil
}
