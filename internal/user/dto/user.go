package dto

import "rewrite/pkg/entity"

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersRequest []UserRequest

func (u *UserRequest) ToEntity() *entity.User {
	return &entity.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type UsersResponse []UserResponse

func (u *UserResponse) FromEntity(entity *entity.User) {
	u.ID = entity.ID
	u.Email = entity.Email
}

func (u *UsersResponse) FromEntity(entities entity.Users) {
	for _, each := range entities {
		var user UserResponse
		user.FromEntity(&each)
		*u = append(*u, user)
	}
}
