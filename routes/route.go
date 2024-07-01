package routes

import (
	"github.com/Lirikku/controllers"
	"github.com/Lirikku/middlewares"
	"github.com/Lirikku/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRoute() *echo.Echo {
	e := echo.New()

	e.Validator = middlewares.NewValidator()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper:        middleware.DefaultSkipper,
		TokenLength:    32,
		TokenLookup:    "header:" + echo.HeaderXCSRFToken,
		ContextKey:     "csrf",
		CookieMaxAge:   86400,
		CookieHTTPOnly: true,
	}))

	// auth
	authGroup := e.Group("/auth")
	authController := controllers.NewAuthController(services.GetAuthRepo())
	{
		authGroup.GET("/register", authController.RegisterView)
		authGroup.GET("/login", authController.LoginView)
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}

	// song lyrics
	songLyricsGroup := e.Group("/song_lyrics", middlewares.Authorized())
	{
		// my song lyrics
		myGroup := songLyricsGroup.Group("/my")
		myController := controllers.NewMySongLyricsController(services.GetMySongLyricsRepo())
		{
			myGroup.Use(middlewares.JWT())
			myGroup.GET("", myController.GetSongLyrics)
			myGroup.GET("/:id", myController.GetSongLyric)
			myGroup.GET("/search", myController.SearchSongLyrics)
			myGroup.POST("", myController.SaveSongLyric)
			myGroup.DELETE("/:id", myController.DeleteSongLyric)
			myGroup.PUT("/:id", myController.UpdateSongLyric)
		}

		// public song lyrics
		publicGroup := songLyricsGroup.Group("/public")
		publicController := controllers.NewPublicSongLyricsController(services.GetPublicSongLyricsRepo())
		{
			publicGroup.GET("/search", publicController.SearchTermSongLyrics)
			publicGroup.POST("/search/audio", publicController.SearchAudioSongLyric)
		}

	}

	return e
}
