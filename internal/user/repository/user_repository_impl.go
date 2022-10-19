package repository

import (
	"context"
	"errors"
	"rewrite/pkg/entity"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExist = errors.New("email already exist")
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (u *UserRepositoryImpl) FindAll(ctx context.Context) (entity.Users, error) {
	var users entity.Users

	err := u.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepositoryImpl) CreateUser(user *entity.User, ctx context.Context) error {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return ErrEmailAlreadyExist
		}

		return err
	}

	return nil
}

func (u *UserRepositoryImpl) FindByEmail(email string, ctx context.Context) (*entity.User, error) {
	var user entity.User

	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
