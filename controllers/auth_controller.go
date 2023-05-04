package controllers

import (
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	service services.IAuthService
}

func NewAuthController(service services.IAuthService) *Auth {
	return &Auth{service}
}

func (a *Auth) Register(c echo.Context) error {
	
	reqAuth := models.UserRegister{}

	c.Bind(&reqAuth)

	if err := c.Validate(&reqAuth); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	err := a.service.CheckUserEmail(reqAuth.Email)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	a.service.CreateUser(reqAuth)
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "success created user",
	})
}

func (a *Auth) Login(c echo.Context) error {
	
	reqAuth := models.UserLogin{}

	c.Bind(&reqAuth)
	
	if err := c.Validate(&reqAuth); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}


	user := models.User{}

	err := configs.DB.First(&user,"email = ?", reqAuth.Email).Error
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "email not registered",
		})
	}

	if !user.CheckPassword(reqAuth.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "invalid email or password",
		})
	}

	token, err := user.GenerateToken()

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success login",
		"token": token,
	})
}