package service

import (
	"context"
	"errors"
	"rewrite/internal/user/dto"
	"rewrite/internal/user/repository"
	"rewrite/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserServiceImpl(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{userRepository}
}

func (u *UserServiceImpl) FindAll(ctx context.Context) (dto.UsersResponse, error) {
	users, err := u.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var dtoUsers dto.UsersResponse
	dtoUsers.FromEntity(users)
	return dtoUsers, nil
}

func (u *UserServiceImpl) CreateUser(user dto.UserRequest, ctx context.Context) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	userEntity := user.ToEntity()

	err = u.userRepository.CreateUser(userEntity, ctx)
	if err != nil {
		if err == repository.ErrEmailAlreadyExist {
			return ErrUserExists
		}
		return err
	}

	return nil
}

func (u *UserServiceImpl) Login(user dto.UserRequest, ctx context.Context) (string, error) {
	userEntity, err := u.userRepository.FindByEmail(user.Email, ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(user.Password))
	if err != nil {
		return "", nil
	}

	token, err := utils.GenerateToken(userEntity)
	if err != nil {
		return "", err
	}

	return token, nil
}
