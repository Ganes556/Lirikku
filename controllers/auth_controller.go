package controllers

import (
	"net/http"

	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	"github.com/Lirikku/utils"
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

	err = a.service.CreateUser(reqAuth)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}
	
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

	user, err := a.service.GetUserByEmail(reqAuth.Email)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "email not registered",
		})
	}

	if !utils.ComparePassword(user.Password,reqAuth.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"message": "incorrect email or password",
		})
	}

	token, err := utils.GenerateToken(user.ID, user.Name)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success login",
		"token": token,
	})
}