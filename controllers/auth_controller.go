package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	"github.com/Lirikku/utils"
	"github.com/Lirikku/view"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	service services.IAuthService
}

func NewAuthController(service services.IAuthService) *Auth {
	return &Auth{service}
}

func (a *Auth) RegisterView(c echo.Context) error {
	return utils.Render(c, http.StatusOK, view.Register(c))
}

func (a *Auth) Register(c echo.Context) error {

	reqAuth := models.UserRegister{}

	c.Bind(&reqAuth)

	if err := c.Validate(&reqAuth); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	err := a.service.CheckUserEmail(reqAuth.Email)

	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": err.Error(),
		})
	}

	err = a.service.CreateUser(reqAuth)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	b, _ := json.Marshal(reqAuth)
	c.Request().Body = io.NopCloser(bytes.NewBuffer(b))
	a.Login(c)

	return c.NoContent(http.StatusNoContent)
}

func (a *Auth) LoginView(c echo.Context) error {
	// csrfToken := c.Get("csrf").(string)
	return utils.Render(c, http.StatusOK, view.Login(c))
}

func (a *Auth) Login(c echo.Context) error {

	reqAuth := models.UserLogin{}

	c.Bind(&reqAuth)

	if err := c.Validate(&reqAuth); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	user, err := a.service.GetUserByEmail(reqAuth.Email)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "incorrect email or password",
		})
	}

	if !user.CheckPassword(reqAuth.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "incorrect email or password",
		})
	}

	store, err := configs.Store.Get(c.Request(), "session")

	if err != nil {
		fmt.Println("err",err)
		return err
	}

	store.Values["auth"] = true
	userData := models.UserJWTDecode{
		ID: user.ID,
		Name: user.Name,
	}
	store.Values["user"] = utils.Convert2Json(userData)
	if err := store.Save(c.Request(), c.Response()); err != nil {
		fmt.Println("err",err)
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusMovedPermanently)
}

func (a *Auth) Logout(c echo.Context) error {

	session, err := configs.Store.Get(c.Request(), "session")

	if err != nil {
		fmt.Println("err",err)
		return err
	}
	
	configs.Store.Delete(c.Request(), c.Response(), session)
	c.Set("auth", nil)
	c.Set("user", nil)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusMovedPermanently)
}