package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"rewrite/internal/user/dto"
	"rewrite/internal/user/service"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindAll(ctx context.Context) (dto.UsersResponse, error) {
	args := m.Called()
	return args.Get(0).(dto.UsersResponse), args.Error(1)
}

func (m *MockUserService) CreateUser(user dto.UserRequest, ctx context.Context) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Login(user dto.UserRequest, ctx context.Context) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

type TestSuiteUserControllers struct {
	suite.Suite
	mockUserService *MockUserService
	userController  *UserController
	echoApp         *echo.Echo
}

func (s *TestSuiteUserControllers) SetupTest() {
	s.mockUserService = new(MockUserService)
	s.userController = NewUserController(s.mockUserService)
	s.echoApp = echo.New()
}

func (s *TestSuiteUserControllers) TearDownTest() {
	s.mockUserService = nil
	s.userController = nil
	s.echoApp = nil
}

func (s *TestSuiteUserControllers) TestInitRoutes() {
	s.NotPanics(func() {
		s.userController.InitRoutes(s.echoApp)
	})
}

func (s *TestSuiteUserControllers) TestGetAllUser() {
	for _, tc := range []struct {
		Name           string
		FunctionUsers  dto.UsersResponse
		FunctionError  error
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name: "Success Get All User",
			FunctionUsers: dto.UsersResponse{
				{
					ID:    1,
					Email: "123@123.com",
				},
				{
					ID:    2,
					Email: "456@456.com",
				},
			},
			ExpectedStatus: 200,
			ExpectedBody: echo.Map{
				"message": "Success getting users",
				"data": []interface{}{
					map[string]interface{}{
						"id":    float64(1),
						"email": "123@123.com",
					},
					map[string]interface{}{
						"id":    float64(2),
						"email": "456@456.com",
					},
				},
			},
		},
		{
			Name:           "Error no user found",
			FunctionUsers:  dto.UsersResponse{},
			ExpectedStatus: 404,
			ExpectedError:  ErrNoUserFound,
		},
		{
			Name:           "Generic error from service",
			FunctionUsers:  dto.UsersResponse{},
			FunctionError:  errors.New("Generic error"),
			ExpectedStatus: 500,
			ExpectedError:  errors.New("Generic error"),
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()
			s.mockUserService.On("FindAll").Return(tc.FunctionUsers, tc.FunctionError)

			r := httptest.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()
			c := s.echoApp.NewContext(r, w)

			err := s.userController.GetAllUser(c)

			if tc.ExpectedError != nil {
				s.Equal(echo.NewHTTPError(tc.ExpectedStatus, tc.ExpectedError.Error()), err)
			} else {
				s.NoError(err)

				var response echo.Map
				err := json.Unmarshal(w.Body.Bytes(), &response)
				s.NoError(err)

				s.Equal(tc.ExpectedStatus, w.Result().StatusCode)
				s.Equal(tc.ExpectedBody, response)
			}

			s.TearDownTest()
		})
	}
}

func (s *TestSuiteUserControllers) TestCreateUser() {
	for _, tc := range []struct {
		Name           string
		RequestBody    interface{}
		RequestContent string
		FunctionError  error
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name: "Success Create User",
			RequestBody: dto.UserRequest{
				Email:    "123@123.com",
				Password: "123",
			},
			RequestContent: "application/json",
			ExpectedStatus: 201,
			ExpectedBody: echo.Map{
				"message": "Success creating user",
			},
		},
		{
			Name: "Error user already exist",
			RequestBody: dto.UserRequest{
				Email:    "123@123.com",
				Password: "123",
			},
			RequestContent: "application/json",
			FunctionError:  service.ErrUserExists,
			ExpectedStatus: 409,
			ExpectedError:  service.ErrUserExists,
		},
		{
			Name:           "Generic error from service",
			RequestBody:    dto.UserRequest{},
			RequestContent: "application/json",
			FunctionError:  errors.New("Generic error"),
			ExpectedStatus: 500,
			ExpectedError:  errors.New("Generic error"),
		},
		{
			Name:           "Error invalid request body",
			RequestBody:    "invalid body",
			RequestContent: "application/json",
			ExpectedStatus: 400,
			ExpectedError:  ErrBadRequestBody,
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
			r.Header.Set("Content-Type", tc.RequestContent)
			w := httptest.NewRecorder()
			c := s.echoApp.NewContext(r, w)

			s.mockUserService.On("CreateUser", tc.RequestBody).Return(tc.FunctionError)
			err = s.userController.CreateUser(c)

			if tc.ExpectedError != nil {
				s.Equal(echo.NewHTTPError(tc.ExpectedStatus, tc.ExpectedError.Error()), err)
			} else {
				s.NoError(err)

				var response echo.Map
				err := json.Unmarshal(w.Body.Bytes(), &response)
				s.NoError(err)

				s.Equal(tc.ExpectedStatus, w.Result().StatusCode)
				s.Equal(tc.ExpectedBody, response)
			}

			s.TearDownTest()
		})
	}
}

func (s *TestSuiteUserControllers) TestLogin() {
	for _, tc := range []struct {
		Name           string
		RequestBody    interface{}
		RequestContent string
		FunctionError  error
		FunctionReturn string
		ExpectedStatus int
		ExpectedBody   echo.Map
		ExpectedError  error
	}{
		{
			Name: "Success Login",
			RequestBody: dto.UserRequest{
				Email:    "123@123.com",
				Password: "123",
			},
			RequestContent: "application/json",
			FunctionReturn: "token",
			ExpectedStatus: 200,
			ExpectedBody: echo.Map{
				"message": "Login success",
				"token":   "token",
			},
		},
		{
			Name: "Error invalid email or password",
			RequestBody: dto.UserRequest{
				Email:    "123@123.com",
				Password: "123",
			},
			RequestContent: "application/json",
			FunctionReturn: "",
			ExpectedStatus: 401,
			ExpectedError:  ErrInvalidCredentials,
		},
		{
			Name:           "Generic error from service",
			RequestBody:    dto.UserRequest{},
			RequestContent: "application/json",
			FunctionError:  errors.New("Generic error"),
			ExpectedStatus: 500,
			ExpectedError:  errors.New("Generic error"),
		},
		{
			Name:           "Error invalid request body",
			RequestBody:    "invalid body",
			RequestContent: "application/json",
			ExpectedStatus: 400,
			ExpectedError:  ErrBadRequestBody,
		},
	} {
		s.Run(tc.Name, func() {
			s.SetupTest()

			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			r := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
			r.Header.Set("Content-Type", tc.RequestContent)
			w := httptest.NewRecorder()
			c := s.echoApp.NewContext(r, w)

			s.mockUserService.On("Login", tc.RequestBody).Return(tc.FunctionReturn, tc.FunctionError)
			err = s.userController.Login(c)

			if tc.ExpectedError != nil {
				s.Equal(echo.NewHTTPError(tc.ExpectedStatus, tc.ExpectedError.Error()), err)
			} else {
				s.NoError(err)

				var response echo.Map
				err := json.Unmarshal(w.Body.Bytes(), &response)
				s.NoError(err)

				s.Equal(tc.ExpectedStatus, w.Result().StatusCode)
				s.Equal(tc.ExpectedBody, response)
			}

			s.TearDownTest()
		})
	}
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(TestSuiteUserControllers))
}
