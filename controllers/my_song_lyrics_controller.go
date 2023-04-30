package controllers

import (
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/Lirikku/utils"
	"github.com/labstack/echo/v4"
)

func GetMySongLyrics(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)

	offset := c.QueryParam("offset")

	offsetInt, err := utils.CheckOffset(offset)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "offset must be a number",
		})
	}

	var resSongLyrics []models.ResponseSongLyric
	
	// load all song lyrics with artist
	configs.DB.Model(&models.SongLyric{}).Limit(5).Offset(offsetInt).Find(&resSongLyrics, "user_id = ?", user.ID)
	
	return c.JSON(http.StatusOK, echo.Map{
		"next": utils.GenerateNextLink(c, offsetInt, len(resSongLyrics)),
		"my_song_lyrics": resSongLyrics,
	})
}

func GetMySongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	var resSongLyric models.ResponseSongLyric
	
	idSongLyric := c.Param("id")

	idSongLyricInt, err := utils.CheckId(idSongLyric)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number",
		})
	}

	// load song lyrics with artist by id song lyrics
	err = configs.DB.Model(&models.SongLyric{}).First(&resSongLyric, "id = ? AND user_id = ?", idSongLyricInt ,user.ID).Error
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": resSongLyric,
	})
}

func SaveMySongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	var newSongLyric models.SongLyric

	c.Bind(&newSongLyric)

	newSongLyric.UserID = user.ID

	err := configs.DB.Create(&newSongLyric).Error

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "song lyric saved successfully",
	})
	
}

func SearchMySongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	offset := c.QueryParam("offset")

	offsetInt, err := utils.CheckOffset(offset)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "offset must be a number",
		})
	}
	
	var resSongLyrics []models.ResponseSongLyric
	
	title := c.QueryParam("title")
	lyric := c.QueryParam("lyric")
	artist_names:= c.QueryParam("artist_names")
	
	configs.DB.Model(&models.SongLyric{}).Where("user_id = ? AND title LIKE ? AND lyric LIKE ? AND artist_names LIKE ?", user.ID, "%"+title+"%", "%"+lyric+"%", "%"+artist_names+"%").Limit(5).Offset(offsetInt).Find(&resSongLyrics)

	return c.JSON(http.StatusOK, echo.Map{
		"next": utils.GenerateNextLink(c, offsetInt, len(resSongLyrics)),
		"my_song_lyrics": resSongLyrics,
	})
}

func DeleteMySongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	idSongLyricInt, err := utils.CheckId(idSongLyric)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number",
		})
	}

	querySongLyric := models.SongLyric{}

	err = configs.DB.First(&querySongLyric, "id = ? AND user_id = ?", idSongLyricInt, user.ID).Error
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	err = configs.DB.Unscoped().Delete(&querySongLyric).Error	

	if err != nil {
		return err
	}
		
	return c.JSON(http.StatusOK, echo.Map{
		"message": "song lyric deleted successfully",
	})

}

func UpdateMySongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	idSongLyricInt, err := utils.CheckId(idSongLyric)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number",
		})
	}
	
	err = configs.DB.First(&models.SongLyric{}, "id = ? AND user_id = ?", idSongLyricInt, user.ID).Error
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	var updateSongLyric models.SongLyric

	c.Bind(&updateSongLyric)

	updateSongLyric.UserID = user.ID

	err = configs.DB.Where("id = ?", idSongLyric).Updates(&updateSongLyric).Error

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "song lyric updated successfully",
	})
}
