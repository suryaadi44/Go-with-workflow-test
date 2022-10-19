package service

import (
	"context"
	"rewrite/internal/user/dto"
)

type UserService interface {
	FindAll(ctx context.Context) (dto.UsersResponse, error)
	CreateUser(user dto.UserRequest, ctx context.Context) error
	Login(user dto.UserRequest, ctx context.Context) (string, error)
}
