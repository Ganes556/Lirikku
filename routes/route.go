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

	e.Static("/", "static")
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
		authGroup.GET("/register", authController.RegisterView).Name = "auth.registerForm"
		authGroup.GET("/login", authController.LoginView).Name = "auth.loginForm"
		authGroup.GET("/logout", authController.Logout).Name = "auth.logout"
		authGroup.POST("/register", authController.Register).Name = "auth.register"
		authGroup.POST("/login", authController.Login).Name = "auth.login"
	}

	publicController := controllers.NewPublicSongLyricsController(services.GetPublicSongLyricsRepo(), services.GetMySongLyricsRepo())
	{
		e.GET("/", publicController.SongLyricsView).Name = "indexSong"
		e.GET("/lyric/:artist/:title", publicController.GetSongDetail).Name = "detailSong"
		e.GET("/search", publicController.SearchSongsByTerm).Name = "searchSong"
		e.POST("/search/audio", publicController.SearchAudioSongLyric).Name = "search.audioSong"
		// e.POST("/search/audiobs64", publicController.SearchBase64SongLyric)
	}

	{
		// my song lyrics
		myGroup := e.Group("/my")
		myController := controllers.NewMySongLyricsController(services.GetMySongLyricsRepo())
		{
			myGroup.GET("", myController.GetSongLyrics).Name = "indexMy"
			myGroup.GET("/:id", myController.GetSongLyric).Name = "my.detailSong"
			myGroup.GET("/search", myController.SearchSongLyrics).Name = "my.searchSong"
			myGroup.POST("", myController.SaveSongLyric).Name = "my.storeSong"
			myGroup.DELETE("/:id", myController.DeleteSongLyric).Name = "my.delSong"
			myGroup.PUT("", myController.UpdateSongLyric).Name = "my.putSong"
			myGroup.POST("/public", publicController.SaveSong).Name = "my.pubSong"
		}

	}

	return e
}
