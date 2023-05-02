package routes

import (
	"github.com/Lirikku/controllers"
	"github.com/Lirikku/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func NewRoute() *echo.Echo{
	e := echo.New()

	e.Validator = middlewares.NewValidatorMiddleware()

	e.Pre(middleware.RemoveTrailingSlash())

	// auth
	authGroup := e.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
	}

	// song lyrics
	songLyricsGroup := e.Group("/song_lyrics")
	{
		// my song lyrics
		mySongLyricsGroup := songLyricsGroup.Group("/my")
		{
			mySongLyricsGroup.Use(middlewares.JwtMiddleware())
			mySongLyricsGroup.GET("", controllers.GetMySongLyrics)
			mySongLyricsGroup.GET("/:id", controllers.GetMySongLyric)
			mySongLyricsGroup.GET("/search", controllers.SearchMySongLyric)
			mySongLyricsGroup.POST("", controllers.SaveMySongLyric)
			mySongLyricsGroup.DELETE("/:id", controllers.DeleteMySongLyric)
			mySongLyricsGroup.PUT("/:id", controllers.UpdateMySongLyric)
		}  

	}

	return e
}