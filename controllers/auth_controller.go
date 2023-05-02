package controllers

import (
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	
	reqAuth := models.UserRegister{}

	c.Bind(&reqAuth)

	if err := c.Validate(&reqAuth); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	err := configs.DB.First(&models.User{},"email = ?", reqAuth.Email).Error
	
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "email already registered",
		})
	}

	newUser := models.User{
		Name: reqAuth.Name,
		Email: reqAuth.Email,
		Password: reqAuth.Password,
	}

	if err := configs.DB.Create(&newUser).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "failed to create user",
		})
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "success created user",
	})
}

func Login(c echo.Context) error {
	
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