package service

import (
	"context"
	"errors"
	"rewrite/internal/user/dto"
	"rewrite/internal/user/repository"
	"rewrite/pkg/entity"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(ctx context.Context) (entity.Users, error) {
	args := m.Called()
	return args.Get(0).(entity.Users), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *entity.User, ctx context.Context) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string, ctx context.Context) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

type TestSuiteUserServices struct {
	suite.Suite
	mockUserRepository *MockUserRepository
	userService        UserService
	ctx                context.Context
}

func (s *TestSuiteUserServices) SetupTest() {
	s.mockUserRepository = new(MockUserRepository)
	s.userService = NewUserServiceImpl(s.mockUserRepository)
	s.ctx = context.Background()
}

func (s *TestSuiteUserServices) TearDownTest() {
	s.mockUserRepository = nil
	s.userService = nil
	s.ctx = nil
}

func (s *TestSuiteUserServices) TestFindAll() {
	for _, tt := range []struct {
		Name           string
		FunctionReturn entity.Users
		FunctionError  error
		ExpectedReturn dto.UsersResponse
		ExpectedErr    error
	}{
		{
			Name: "Success",
			FunctionReturn: entity.Users{
				entity.User{
					Email:    "123@123.com",
					Password: "123",
				},
				entity.User{
					Email:    "456@456.com",
					Password: "456",
				},
			},
			FunctionError: nil,
			ExpectedReturn: dto.UsersResponse{
				dto.UserResponse{
					Email: "123@123.com",
				},
				dto.UserResponse{
					Email: "456@456.com",
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:           "Generic Error from Repository",
			FunctionReturn: entity.Users{},
			FunctionError:  errors.New("Generic Error"),
			ExpectedReturn: nil,
			ExpectedErr:    errors.New("Generic Error"),
		},
	} {
		s.SetupTest()
		s.Run(tt.Name, func() {
			s.mockUserRepository.On("FindAll").Return(tt.FunctionReturn, tt.FunctionError)
			result, err := s.userService.FindAll(s.ctx)
			s.Equal(tt.ExpectedReturn, result)
			s.Equal(tt.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserServices) TestCreateUser() {
	for _, tt := range []struct {
		Name          string
		FunctionError error
		UserRequest   dto.UserRequest
		ExpectedErr   error
	}{
		{
			Name:          "Success",
			FunctionError: nil,
			UserRequest: dto.UserRequest{
				Email:    "123@13.com",
				Password: "123",
			},
			ExpectedErr: nil,
		},
		{
			Name:          "User email already exists",
			FunctionError: repository.ErrEmailAlreadyExist,
			UserRequest:   dto.UserRequest{},
			ExpectedErr:   ErrUserExists,
		},
		{
			Name:          "Generic Error from Repository",
			FunctionError: errors.New("Generic Error"),
			UserRequest:   dto.UserRequest{},
			ExpectedErr:   errors.New("Generic Error"),
		},
	} {
		s.SetupTest()
		s.Run(tt.Name, func() {
			s.mockUserRepository.On("CreateUser", mock.Anything).Return(tt.FunctionError)
			err := s.userService.CreateUser(tt.UserRequest, s.ctx)
			s.Equal(tt.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserServices) TestLogin() {
	for _, tt := range []struct {
		Name           string
		FunctionReturn *entity.User
		FunctionError  error
		UserRequest    dto.UserRequest
		ExpectedErr    error
	}{
		{
			Name: "Success",
			FunctionReturn: &entity.User{
				Email:    "123@123.com",
				Password: "123",
			},
			FunctionError: nil,
			UserRequest: dto.UserRequest{
				Email:    "123@123.com",
				Password: "123",
			},
			ExpectedErr: nil,
		},
		{
			Name:           "User not found",
			FunctionReturn: nil,
			FunctionError:  gorm.ErrRecordNotFound,
			UserRequest:    dto.UserRequest{},
		},
		{
			Name:           "Generic Error from Repository",
			FunctionReturn: nil,
			FunctionError:  errors.New("Generic Error"),
			UserRequest:    dto.UserRequest{},
			ExpectedErr:    errors.New("Generic Error"),
		},
	} {
		s.SetupTest()
		s.Run(tt.Name, func() {
			s.mockUserRepository.On("FindByEmail", mock.Anything).Return(tt.FunctionReturn, tt.FunctionError)
			_, err := s.userService.Login(tt.UserRequest, s.ctx)
			s.Equal(tt.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestSuiteUserServices))
}
