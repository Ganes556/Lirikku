package routes

import (
	"github.com/Lirikku/controllers"
	"github.com/Lirikku/middlewares"
	"github.com/Lirikku/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func NewRoute() *echo.Echo{
	e := echo.New()

	e.Validator = middlewares.NewValidator()

	e.Pre(middleware.RemoveTrailingSlash())
	
	// auth
	authGroup := e.Group("/auth")
	authController := controllers.NewAuthController(services.GetAuthRepo())
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}
	
	// song lyrics
	songLyricsGroup := e.Group("/song_lyrics")
	{
		// my song lyrics
		mySongLyricsGroup := songLyricsGroup.Group("/my")
		mySongLyricsController := controllers.NewMySongLyricsController(services.GetMySongLyricsRepo())
		{
			mySongLyricsGroup.Use(middlewares.JWT())
			mySongLyricsGroup.GET("", mySongLyricsController.GetAll)
			mySongLyricsGroup.GET("/:id", mySongLyricsController.Get)
			mySongLyricsGroup.GET("/search", mySongLyricsController.Search)
			mySongLyricsGroup.POST("", mySongLyricsController.Save)
			mySongLyricsGroup.DELETE("/:id", mySongLyricsController.Delete)
			mySongLyricsGroup.PUT("/:id", mySongLyricsController.Update)
		}  

		// public song lyrics
		publicSongLyricsGroup := songLyricsGroup.Group("/public")
		publicSongLyricsController := controllers.NewPublicSongLyricsController(services.GetPublicSongLyricsRepo())
		{
			publicSongLyricsGroup.GET("/search", publicSongLyricsController.SearchTerm)			
			publicSongLyricsGroup.POST("/search/audio", publicSongLyricsController.SearchAudio)
		}

	}

	return e
}