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

	e.HTTPErrorHandler = controllers.CustomHTTPErrorHandler
	e.Validator = middlewares.NewValidator()

	e.Pre(middleware.RemoveTrailingSlash())

	// csrf middleware
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper:        middleware.DefaultSkipper,
		TokenLength:    32,
		TokenLookup:    "header:" + echo.HeaderXCSRFToken,
		ContextKey:     "csrf",
		CookieMaxAge:   86400,
		CookieHTTPOnly: true,
	}))

	// auth middleware
	e.Use(middlewares.Authorized())

	// auth
	authGroup := e.Group("/auth")
	authController := controllers.NewAuthController(services.GetAuthRepo())
	{
		authGroup.GET("/register", authController.RegisterView)
		authGroup.GET("/login", authController.LoginView)
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}

	publicController := controllers.NewPublicSongLyricsController(services.GetPublicSongLyricsRepo())
	{
		e.GET("/", publicController.SongLyricsView)
		e.GET("/lyric/:artist/:title", publicController.GetSongDetail)
		e.GET("/search", publicController.SearchSongsByTerm)
		e.POST("/search/audio", publicController.SearchAudioSongLyric)
		e.POST("/search/audiobs64", publicController.SearchBase64SongLyric)
	}

	{
		// my song lyrics
		myGroup := e.Group("/my")
		myController := controllers.NewMySongLyricsController(services.GetMySongLyricsRepo())
		{
			myGroup.GET("", myController.GetSongLyrics)
			myGroup.GET("/:id", myController.GetSongLyric)
			myGroup.GET("/search", myController.SearchSongLyrics)
			myGroup.POST("", myController.SaveSongLyric)
			myGroup.DELETE("/:id", myController.DeleteSongLyric)
			myGroup.PUT("/:id", myController.UpdateSongLyric)
		}

	}

	return e
}
