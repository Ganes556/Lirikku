package routes

import (
	"github.com/Lirikku/controllers"
	"github.com/Lirikku/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func NewRoute() *echo.Echo{
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	// auth
	authGroup := e.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
	}

	// song_lyric
	songLyricGroup := e.Group("/song_lyrics")
	{
		songLyricGroup.Use(middlewares.JwtMiddleware())
	}
	
	return e
}