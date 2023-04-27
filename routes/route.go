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
	mySongLyricGroup := e.Group("/my_song_lyrics")
	{
		mySongLyricGroup.Use(middlewares.JwtMiddleware())
		mySongLyricGroup.GET("", controllers.GetMySongLyrics)
		mySongLyricGroup.POST("", controllers.SaveMySongLyric)
		mySongLyricGroup.DELETE("/:id", controllers.DeleteMySongLyric)
		mySongLyricGroup.GET("/search", controllers.SearchMySongLyric)
	}   
	
	// globaly song lyric
	
	// songLyricGroup := e.Group("/song_lyrics")
	// {
	// 	songLyricGroup.Use(middlewares.JwtMiddleware())
	// 	songLyricGroup.GET("/", controllers.GetMySongLyrics)
	// 	songLyricGroup.POST("/", controllers.SaveMySongLyric)
	// }



	return e
}