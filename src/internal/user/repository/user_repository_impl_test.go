package repository

import (
	"context"
	"errors"
	"regexp"
	"rewrite/pkg/entity"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestSuiteUserRepository struct {
	suite.Suite
	Mock           sqlmock.Sqlmock
	userRepository UserRepository
	ctx            context.Context
}

func (s *TestSuiteUserRepository) SetupTest() {
	dbMock, mock, err := sqlmock.New()
	s.NoError(err)

	DB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      dbMock,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	s.NoError(err)

	s.Mock = mock
	s.userRepository = NewUserRepositoryImpl(DB)
	s.ctx = context.Background()
}

func (s *TestSuiteUserRepository) TeardownTest() {
	s.Mock = nil
	s.userRepository = nil
	s.ctx = nil
}

func (s *TestSuiteUserRepository) TestFindAll() {
	for _, tt := range []struct {
		Name           string
		Query          string
		Rows           *sqlmock.Rows
		Err            error
		ExpectedReturn entity.Users
		ExpectedErr    error
	}{
		{
			Name:  "Success",
			Query: "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL",
			Rows: sqlmock.NewRows([]string{"email", "password"}).
				AddRow("123@123.com", "123").
				AddRow("456@456.com", "456"),
			ExpectedReturn: entity.Users{
				{
					Email:    "123@123.com",
					Password: "123",
				},
				{
					Email:    "456@456.com",
					Password: "456",
				},
			},
		},
		{
			Name:           "Generic Error from DB",
			Query:          "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL",
			Rows:           nil,
			Err:            errors.New("generic error"),
			ExpectedReturn: nil,
			ExpectedErr:    errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tt.Name, func() {
			if tt.Err != nil {
				s.Mock.ExpectQuery(regexp.QuoteMeta(tt.Query)).WillReturnError(tt.Err)
			} else {
				s.Mock.ExpectQuery(regexp.QuoteMeta(tt.Query)).WillReturnRows(tt.Rows)
			}

			result, err := s.userRepository.FindAll(s.ctx)

			s.Equal(tt.ExpectedReturn, result)
			s.Equal(tt.ExpectedErr, err)
		})
		s.TeardownTest()
	}
}

func (s *TestSuiteUserRepository) TestCreateUser() {
	for _, tt := range []struct {
		Name        string
		Query       string
		Err         error
		ExpectedErr error
	}{
		{
			Name:  "Success",
			Query: "INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`email`,`password`) VALUES (?,?,?,?,?)",
		},
		{
			Name:        "Generic Error from DB",
			Query:       "INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`email`,`password`) VALUES (?,?,?,?,?)",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tt.Name, func() {
			s.Mock.ExpectBegin()
			if tt.Err != nil {
				s.Mock.ExpectExec(regexp.QuoteMeta(tt.Query)).WillReturnError(tt.Err)
				s.Mock.ExpectRollback()
			} else {
				s.Mock.ExpectExec(regexp.QuoteMeta(tt.Query)).WillReturnResult(sqlmock.NewResult(1, 1))
				s.Mock.ExpectCommit()
			}

			err := s.userRepository.CreateUser(&entity.User{
				Model:    gorm.Model{},
				Email:    "",
				Password: "",
			}, s.ctx)

			s.Equal(tt.ExpectedErr, err)
		})
		s.TeardownTest()
	}
}

func (s *TestSuiteUserRepository) TestFindByEmail() {
	for _, tt := range []struct {
		Name           string
		Email          string
		Query          string
		Rows           *sqlmock.Rows
		Err            error
		ExpectedReturn *entity.User
		ExpectedErr    error
	}{
		{
			Name:  "Success",
			Email: "123@123.com",
			Query: "SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1",
			Rows: sqlmock.NewRows([]string{"email", "password"}).
				AddRow("123@123.com", "123"),
			ExpectedReturn: &entity.User{
				Email:    "123@123.com",
				Password: "123",
			},
		},
		{
			Name:           "Generic Error from DB",
			Email:          "123@123.com",
			Query:          "SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1",
			Rows:           nil,
			Err:            errors.New("generic error"),
			ExpectedReturn: nil,
			ExpectedErr:    errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tt.Name, func() {
			if tt.Err != nil {
				s.Mock.ExpectQuery(regexp.QuoteMeta(tt.Query)).WillReturnError(tt.Err)
			} else {
				s.Mock.ExpectQuery(regexp.QuoteMeta(tt.Query)).WillReturnRows(tt.Rows)
			}

			result, err := s.userRepository.FindByEmail(tt.Email, s.ctx)

			s.Equal(tt.ExpectedReturn, result)
			s.Equal(tt.ExpectedErr, err)
		})
		s.TeardownTest()
	}
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteUserRepository))
}
