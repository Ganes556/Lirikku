package routes

import (
	"github.com/Lirikku/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func NewRoute() *echo.Echo{
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	// auth
	authGroup := e.Group("/auth")
	authGroup.POST("/register", controllers.Register)
	authGroup.POST("/login", controllers.Login)

	return e
}