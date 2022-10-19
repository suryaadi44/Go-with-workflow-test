package dto

import (
	"rewrite/pkg/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRequest_ToEntity(t *testing.T) {
	tests := []struct {
		name string
		u    *UserRequest
		want *entity.User
	}{
		{
			name: "UserRequest ToEntity",
			u: &UserRequest{
				Email:    "123@123.com",
				Password: "123",
			},
			want: &entity.User{
				Email:    "123@123.com",
				Password: "123",
			},
		},
		{
			name: "UserRequest ToEntity with empty string",
			u: &UserRequest{
				Email:    "",
				Password: "",
			},
			want: &entity.User{
				Email:    "",
				Password: "",
			},
		},
		{
			name: "UserRequest ToEntity with empty field",
			u:    &UserRequest{},
			want: &entity.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.u.ToEntity())
		})
	}
}

func TestUserResponse_FromEntity(t *testing.T) {
	tests := []struct {
		name   string
		want   *UserResponse
		entity *entity.User
	}{
		{
			name: "UserResponse FromEntity",
			want: &UserResponse{
				Email: "123@123.com",
			},
			entity: &entity.User{
				Email: "123@123.com",
			},
		},
		{
			name: "UserResponse FromEntity with empty string",
			want: &UserResponse{
				Email: "",
			},
			entity: &entity.User{
				Email: "",
			},
		},
		{
			name:   "UserResponse FromEntity with empty field",
			want:   &UserResponse{},
			entity: &entity.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserResponse{}
			u.FromEntity(tt.entity)

			assert.Equal(t, tt.want, u)
		})
	}
}

func TestUsersResponse_FromEntity(t *testing.T) {
	tests := []struct {
		name   string
		want   *UsersResponse
		entity entity.Users
	}{
		{
			name: "UsersResponse FromEntity",
			want: &UsersResponse{
				{
					Email: "123@123.com",
				},
				{
					Email: "456@456.com",
				},
			},
			entity: entity.Users{
				{
					Email: "123@123.com",
				},
				{
					Email: "456@456.com",
				},
			},
		},
		{
			name: "UsersResponse FromEntity with empty string",
			want: &UsersResponse{
				{
					Email: "",
				},
				{
					Email: "",
				},
			},
			entity: entity.Users{
				{
					Email: "",
				},
				{
					Email: "",
				},
			},
		},
		{
			name:   "UsersResponse FromEntity with empty field",
			want:   &UsersResponse{},
			entity: entity.Users{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UsersResponse{}
			u.FromEntity(tt.entity)

			assert.Equal(t, tt.want, u)
		})
	}
}
