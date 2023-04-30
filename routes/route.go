package routes

import (
	"github.com/Lirikku/controllers"
	"github.com/labstack/echo/v4"
)


func NewRoute() *echo.Echo{
	e := echo.New()

	// auth
	authGroup := e.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
	}

	return e
}