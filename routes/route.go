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

	// my song lyrics
	songLyricGroup := e.Group("/song_lyrics")
	{
		mySongLyricGroup := songLyricGroup.Group("/my")
		{
			mySongLyricGroup.Use(middlewares.JwtMiddleware())
			mySongLyricGroup.GET("", controllers.GetMySongLyrics)
			mySongLyricGroup.GET("/:id", controllers.GetMySongLyric)
			mySongLyricGroup.GET("/search", controllers.SearchMySongLyric)
			mySongLyricGroup.POST("", controllers.SaveMySongLyric)
			mySongLyricGroup.DELETE("/:id", controllers.DeleteMySongLyric)
			mySongLyricGroup.PUT("/:id", controllers.UpdateMySongLyric)
		}
	}   

	return e
}