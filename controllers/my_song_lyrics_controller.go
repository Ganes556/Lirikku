package controllers

import (
	"net/http"
	"net/url"
	"strconv"

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

func (my *MySongLyrics) GetAll(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)

	offset := c.QueryParam("offset")

	offsetInt, err := utils.CheckOffset(offset)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "offset must be a number",
		})
	}

	resSongLyrics, err := my.service.GetMySongLyrics(user.ID,offsetInt)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "failed to get data",
		})
	}
	
	next := utils.GenerateNextLink(c, len(resSongLyrics), url.Values{
		"offset": {strconv.Itoa(offsetInt + 5)},
	}.Encode())

	return c.JSON(http.StatusOK, echo.Map{
		"next": next,
		"my_song_lyrics": resSongLyrics,
	})
}

func (my *MySongLyrics) Get(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")
	
	idSongLyricInt, err := utils.CheckId(idSongLyric)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number",
		})
	}

	resSongLyric, err := my.service.GetMySongLyric(idSongLyricInt, user.ID)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"my_song_lyrics": resSongLyric,
	})
}

func (my *MySongLyrics) Save(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	var reqSongLyricWrite models.SongLyricWrite

	c.Bind(&reqSongLyricWrite)

	my.service.SaveMySongLyric(user.ID, reqSongLyricWrite)

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "song lyric saved successfully",
	})
	
}

func (my *MySongLyrics) Search(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	offset := c.QueryParam("offset")

	offsetInt, err := utils.CheckOffset(offset)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "offset must be a number",
		})
	}
	
	
	title := c.QueryParam("title")
	lyric := c.QueryParam("lyric")
	artist_names:= c.QueryParam("artist_names")

	resSongLyrics, _ := my.service.SearchMySongLyrics(user.ID, title, lyric, artist_names, offsetInt)

	next := utils.GenerateNextLink(c, len(resSongLyrics), url.Values{
		"title": {title},
		"lyric": {lyric},
		"artist_names": {artist_names},
		"offset": {strconv.Itoa(offsetInt + 5)},
	}.Encode())

	return c.JSON(http.StatusOK, echo.Map{
		"next": next,
		"my_song_lyrics": resSongLyrics,
	})
}

func (my *MySongLyrics) Delete(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	idSongLyricInt, err := utils.CheckId(idSongLyric)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number",
		})
	}

	_, err = my.service.GetMySongLyric(idSongLyricInt, user.ID)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	my.service.DeleteMySongLyric(idSongLyricInt, user.ID)
		
	return c.JSON(http.StatusOK, echo.Map{
		"message": "song lyric deleted successfully",
	})

}

func (my *MySongLyrics) Update(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)
	
	idSongLyric := c.Param("id")

	idSongLyricInt, err := utils.CheckId(idSongLyric)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "id must be a number",
		})
	}
	
	_, err = my.service.GetMySongLyric(idSongLyricInt, user.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "song lyric not found",
		})
	}

	var reqSongLyricWrite models.SongLyricWrite

	c.Bind(&reqSongLyricWrite)

	my.service.UpdateMySongLyric(idSongLyricInt, user.ID, reqSongLyricWrite)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "song lyric updated successfully",
	})
}
