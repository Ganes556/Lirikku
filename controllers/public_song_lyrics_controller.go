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

func (pub *PublicSongLyrics) SearchTerm(c echo.Context) error {
	term := c.QueryParam("term")
	offset := c.QueryParam("offset")
	
	offsetInt, err := utils.CheckOffset(offset)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "offset must be a number",
		})
	}

	resPublicSongLyrics, _ := pub.service.SearchByTerm(term, "artists,songs", "5", offsetInt)
	
	next := utils.GenerateNextLink(c, len(resPublicSongLyrics), url.Values{
		"term": {term},
		"offset": {strconv.Itoa(offsetInt + 5)},
	}.Encode())

	return c.JSON(http.StatusOK, echo.Map{
		"next": next,
		"public_song_lyrics": resPublicSongLyrics,
	})

}

func (pub *PublicSongLyrics) SearchAudio(c echo.Context) error {
	audioData, err := c.FormFile("audio")

	if err != nil {
		return err
	}

	if audioData.Size > 500000 {
		return c.JSON(http.StatusRequestEntityTooLarge, echo.Map{
			"message": "audio size must be less than 500kb",
		})
	}

	isAudio, err := utils.CheckAudioFile(audioData)

	if err != nil {
		return err
	}

	if !isAudio {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid file type. please upload an audio file.",
		})
	}
	
	rawBases64, err := utils.Audio2RawBase64(audioData)

	if err != nil {
		return err
	}

	resData, err := pub.service.SearchByAudio(rawBases64)
	
	if err != nil {
		if err.Error() == "song not found" {
			return echo.NewHTTPError(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"public_song_lyrics": resData,
	})
}
