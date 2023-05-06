package mocks

import (
	"github.com/Lirikku/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) CreateUser(req models.UserRegister) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthRepo) CheckUserEmail(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthRepo) GetUserByEmail(email string) (models.User, error){
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}	

