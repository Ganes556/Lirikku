package controllers

import (
	"net/http"

	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	"github.com/Lirikku/utils"
	"github.com/labstack/echo/v4"
)

type MySongLyrics struct {
	service services.IMySongLyricsService
}

func NewMySongLyricsController(service services.IMySongLyricsService) *MySongLyrics {
	return &MySongLyrics{service}
}

func (my *MySongLyrics) GetSongLyrics(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	pageSize, offset := utils.GetPageSizeAndOffset(c)
	
	resSongLyrics, err := my.service.GetSongLyrics(user.ID, offset, pageSize)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	
	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": resSongLyrics,
	})
}

func (my *MySongLyrics) GetSongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")
	
	idSongLyricInt := utils.CheckId(idSongLyric)
	
	if idSongLyricInt == -1{
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number and greater than 0",
		})
	}

	resSongLyric, err := my.service.GetSongLyric(idSongLyricInt, user.ID)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{
			"message": "song lyric not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": resSongLyric,
	})
}

func (my *MySongLyrics) SaveSongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	var reqSongLyricWrite models.SongLyricWrite
	
	c.Bind(&reqSongLyricWrite)
	
	if err := c.Validate(reqSongLyricWrite); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	err := my.service.CheckSongLyric(user.ID, reqSongLyricWrite)

	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, echo.Map{
			"message": err.Error(),
		})
	}

	err = my.service.SaveSongLyric(user.ID, reqSongLyricWrite)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "song lyric saved successfully",
	})
	
}

func (my *MySongLyrics) SearchSongLyrics(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	pageSize, offset := utils.GetPageSizeAndOffset(c)
	
	title := c.QueryParam("title")
	lyric := c.QueryParam("lyric")
	artist_names:= c.QueryParam("artist_names")
	
	resSongLyrics, err := my.service.SearchSongLyrics(user.ID, title, lyric, artist_names, offset, pageSize)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	if len(resSongLyrics) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{
			"message": "song lyric not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": resSongLyrics,
	})
}

func (my *MySongLyrics) DeleteSongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	idSongLyricInt := utils.CheckId(idSongLyric)

	if idSongLyricInt == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number and greater than 0",
		})
	}

	_, err := my.service.GetSongLyric(idSongLyricInt, user.ID)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{
			"message": "song lyric not found",
		})
	}

	err = my.service.DeleteSongLyric(idSongLyricInt, user.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}
		
	return c.JSON(http.StatusOK, echo.Map{
		"message": "song lyric deleted successfully",
	})

}

func (my *MySongLyrics) UpdateSongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	idSongLyricInt := utils.CheckId(idSongLyric)

	if idSongLyricInt == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number and greater than 0",
		})
	}
	
	_, err := my.service.GetSongLyric(idSongLyricInt, user.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{
			"message": "song lyric not found",
		})
	}

	var reqSongLyricWrite models.SongLyricWrite

	c.Bind(&reqSongLyricWrite)

	err = my.service.UpdateSongLyric(idSongLyricInt, user.ID, reqSongLyricWrite)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "song lyric updated successfully",
	})
}
