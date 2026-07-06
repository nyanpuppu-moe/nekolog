package repository

import (
	"nekolog/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	dataBase *gorm.DB
}

func NewUserRepository(dataBase *gorm.DB) *UserRepository {
	return &UserRepository{
		dataBase: dataBase,
	}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.dataBase.
		Create(user).
		Error
}

func (r *UserRepository) FindByID(id model.UserID) (*model.User, error) {
	var user model.User

	err := r.dataBase.
		Where(&model.User{
			ID: id,
		}).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByName(name string) (*model.User, error) {
	var user model.User

	err := r.dataBase.
		Where(&model.User{
			Name: name,
		}).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
