package services

import (
	"errors"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
)

type IAuthService interface {
	CheckUserEmail(email string) error
	CreateUser(req models.UserRegister) error
	GetUserByEmail(email string) (models.User, error)
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

func (ar *AuthRepo) CheckUserEmail(email string) error {

	err := configs.DB.First(&models.User{}, "email = ?", email).Error
	if err == nil {
		return errors.New("email already registered")
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

func (ar *AuthRepo) GetUserByEmail(email string) (models.User, error){
	var user models.User

	err := configs.DB.First(&user, "email = ?", email ).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

