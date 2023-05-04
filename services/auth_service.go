package services

import (
	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
)

type IAuthService interface {
	GetUserByEmail(email string) error
	CreateUser(req models.UserRegister) error
}

type AuthRepo struct {}

var authRepo IAuthService


func init() {
	authRepo = &AuthRepo{}
}

func GetAuthRepo() IAuthService {
	return authRepo
}

func SetAuthRepo(repo IAuthService) {
	authRepo = repo
}

func (ar *AuthRepo) GetUserByEmail(email string) error {
	err := configs.DB.First(&models.User{}, "email = ?", email).Error
	
	if err != nil {
		return err
	}

	return nil
}

func (ar *AuthRepo) CreateUser(req models.UserRegister) error {
	newUser := models.User{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
	}

	if err := configs.DB.Create(&newUser).Error; err != nil {
		return err
	}

	return nil
}
