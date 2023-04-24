package controllers

import (
	"fmt"
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func GetMySongLyrics(c echo.Context) error {
	data := c.Get("user").(*jwt.Token).Claims.(*models.JWTClaims)
	var songLyrics []models.SongLyric
	
	configs.DB.Model(songLyrics).Where("user_id = ?", data.ID).Find(&songLyrics)
	
	return c.JSON(200, echo.Map{
		"my_song_lyrics": songLyrics,
	})
}

func SaveMySongLyric(c echo.Context) error {
	data := c.Get("user").(*jwt.Token).Claims.(*models.JWTClaims)
	
	var reqSongLyric models.ReqSongLyric

	c.Bind(&reqSongLyric)

	var songLyric models.SongLyric

	songLyric.Title = reqSongLyric.Title
	

	// songLyric

	// configs.DB.Create(&songLyric)

	// err := configs.DB.Create(&songLyric).Error

	// if err != nil {
	// 	return err
	// }

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Song lyric saved successfully",
	})
	
}