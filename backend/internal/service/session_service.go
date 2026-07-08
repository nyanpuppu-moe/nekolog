package service

import (
	"errors"
	"strings"

	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"nekolog/internal/dto"
	"nekolog/internal/model"
	"nekolog/internal/repository"

	"golang.org/x/crypto/argon2"
)

type SessionService struct {
	userRepo *repository.UserRepository
}

func NewSessionService(repo *repository.UserRepository) *SessionService {
	return &SessionService{userRepo: repo}
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		1,       // time
		64*1024, // memory (64MB)
		4,       // threads
		32,      // key length
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return b64Salt + "." + b64Hash, nil
}

func comparePassword(password, encoded string) bool {
	parts := strings.Split(encoded, ".")

	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		1,
		64*1024,
		4,
		32,
	)

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
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

	hashedPassword, err := hashPassword(req.Password)
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

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if comparePassword(user.Password, hashedPassword) {
		return nil, errors.New("Incorrect email or password")
	}

	return user, nil
}
