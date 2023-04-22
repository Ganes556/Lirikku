package controllers

import (
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	
	reqAuth := models.ReqAuthUser{}

	c.Bind(&reqAuth)

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
