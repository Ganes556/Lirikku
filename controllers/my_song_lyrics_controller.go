package controllers

import (
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetMySongLyrics(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	var songLyrics []models.SongLyric
	
	// load all song lyrics with artist
	configs.DB.Preload("Artists", func(tx *gorm.DB) *gorm.DB{
		return tx.Omit("created_at,deleted_at,updated_at")
	}).Omit("created_at,deleted_at,updated_at").Where("user_id = ?", user.ID).Find(&songLyrics)

	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": songLyrics,
	})
}

func GetMySongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	var songLyrics []models.SongLyric
	
	idSongLyric := c.Param("id")

	// load song lyrics with artist by id song lyrics
	configs.DB.
	Preload("Artists", func(tx *gorm.DB) *gorm.DB{
		return tx.Omit("created_at,deleted_at,updated_at")
	}).
	Omit("created_at,deleted_at,updated_at").
	Where("id = ? AND user_id = ?", idSongLyric, user.ID).
	Find(&songLyrics)

	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": songLyrics,
	})
}

func SaveMySongLyric(c echo.Context) error {
	user := c.Get("user").(models.UserJWTDecode)
	
	var songLyric models.WriteSongLyric

	c.Bind(&songLyric)

	var newSongLyric = songLyric.Convert2SongLyric(user.ID)

	err := configs.DB.Create(newSongLyric).Error

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Song lyric saved successfully",
	})
	
}

func SearchMySongLyric(c echo.Context) error {
	user := c.Get("user").(models.UserJWTDecode)
	
	
	var songLyrics []models.SongLyric
	
	title := c.QueryParam("title")
	lyric := c.QueryParam("lyric")
	artists:= c.QueryParam("artists")	

	querySql := `SELECT sl.*, a.* FROM song_lyric_artists sla 
					JOIN song_lyrics sl ON sl.id = sla.song_lyric_id 
					JOIN artists a ON a.id = sla.artist_id
					WHERE sl.user_id = ? 
					AND sl.title LIKE ?
					AND sl.lyric LIKE ?
					AND a.name LIKE ?`

	configs.DB.Raw(querySql, user.ID, "%"+title+"%", "%"+lyric+"%","%"+artists+"%").Preload("Artists").Find(&songLyrics)
	
	// load all song lyrics with artist
	// configs.DB.Debug().Preload("Artists").Where("user_id = ? AND (title LIKE ? OR lyric LIKE ? OR artists.name LIKE ?)", user.ID, "%"+title+"%", "%"+lyric+"%","%"+artists+"%").Find(&songLyrics)

	// configs.DB.Debug().Raw("SELECT sl.*, a.* FROM song_lyric_artists sla JOIN song_lyrics sl ON sl.id = sla.song_lyric_id JOIN artists a ON a.id = sla.artist_id WHERE sl.user_id = ?;", user.ID).Scan(&songLyrics)
	
	// configs.DB.Debug().Table("song_lyric_artists").Joins("Artists").Joins("SongLyrics").Find(&songLyrics)

	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": songLyrics,
	})
}

func DeleteMySongLyric(c echo.Context) error {
	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	err := configs.DB.First(&models.SongLyric{}, "id = ? AND user_id = ?", idSongLyric, user.ID).Error
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	var querySongLyric *models.SongLyric
	
	configs.DB.Where("id = ? AND user_id = ?", idSongLyric, user.ID).Preload("Artists").First(&querySongLyric)

	err = configs.DB.Unscoped().Delete(&querySongLyric).Error	

	if err != nil {
		return err
	}
		
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Song lyric deleted successfully",
	})

}

func UpdateMySongLyric(c echo.Context) error {
	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	err := configs.DB.First(&models.SongLyric{}, "id = ? AND user_id = ?", idSongLyric, user.ID).Error
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	var updateSongLyric models.WriteSongLyric

	c.Bind(&updateSongLyric)

	newUpdatedSongLyric := updateSongLyric.Convert2SongLyric(user.ID)

	err = configs.DB.Where("id = ?", idSongLyric).Updates(&newUpdatedSongLyric).Error

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Song lyric updated successfully",
	})
}
