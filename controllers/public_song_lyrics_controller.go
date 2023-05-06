package controllers

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/Lirikku/services"
	"github.com/Lirikku/utils"
	"github.com/labstack/echo/v4"
)

type PublicSongLyrics struct {
	service services.IPublicSongLyricsService
}

func NewPublicSongLyricsController(service services.IPublicSongLyricsService) *PublicSongLyrics {
	return &PublicSongLyrics{service}
}

func (pub *PublicSongLyrics) SearchTermSongLyrics(c echo.Context) error {
	term := c.QueryParam("term")
	offset := c.QueryParam("offset")
	
	offsetInt, err := utils.CheckOffset(offset)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "offset must be a number",
		})
	}

	resPublicSongLyrics, err := pub.service.SearchSongLyricsByTermShazam(term, "artists,songs", "5", offsetInt)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	next := utils.GenerateNextLink(c, len(resPublicSongLyrics), url.Values{
		"term": {term},
		"offset": {strconv.Itoa(offsetInt + 5)},
	}.Encode())

	return c.JSON(http.StatusOK, echo.Map{
		"next": next,
		"public_song_lyrics": resPublicSongLyrics,
	})

}

func (pub *PublicSongLyrics) SearchAudioSongLyric(c echo.Context) error {
	audioData, _ := c.FormFile("audio")

	isAudio := utils.CheckAudioFile(audioData)
	if !isAudio {
		// log.Println(isAudio)
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "invalid file type. please upload an audio file",
		})
	}
	
	if audioData.Size > 500000 {
		return echo.NewHTTPError(http.StatusRequestEntityTooLarge, echo.Map{
			"message": "audio size must be less than 500kb",
		})
	}

	rawBases64 := utils.Audio2RawBase64(audioData)

	resData, err := pub.service.SearchSongLyricByAudioRapidShazam(rawBases64)
	
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{
			"message": "song lyric not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"public_song_lyrics": resData,
	})
}
