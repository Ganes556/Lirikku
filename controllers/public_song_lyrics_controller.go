package controllers

import (
	"net/http"

	"github.com/Lirikku/services"
	"github.com/Lirikku/utils"
	"github.com/Lirikku/view"
	"github.com/labstack/echo/v4"
)

type PublicSongLyrics struct {
	service services.IPublicSongLyricsService
}

func NewPublicSongLyricsController(service services.IPublicSongLyricsService) *PublicSongLyrics {
	return &PublicSongLyrics{service}
}

func (pub *PublicSongLyrics) SongLyricsView(c echo.Context) error {
	return utils.Render(c, http.StatusOK, view.SongLyric())
}

func (pub *PublicSongLyrics) SearchSongsByTerm(c echo.Context) error {
	term := c.QueryParam("term")

	pageSize, offset := utils.GetPageSizeAndOffset(c)

	res, err := pub.service.SearchSongsByTermShazam(term, "artists,songs", offset, pageSize)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	if len(c.Request().Header.Values("HX-Request")) > 0 && c.Request().Header.Values("HX-Request")[0] == "true" {
		return utils.Render(c, http.StatusOK, view.ResultSearch(res))
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})

}

func (pub *PublicSongLyrics) GetSongDetail(c echo.Context) error {
	artist := c.Param("artist")
	title := c.Param("title")

	res, err := pub.service.GetSongDetail(artist, title)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	if len(c.Request().Header.Values("HX-Request")) > 0 && c.Request().Header.Values("HX-Request")[0] == "true" {
		return utils.Render(c, http.StatusOK, view.SongsDetail(res))
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})

}

func (pub *PublicSongLyrics) SearchAudioSongLyric(c echo.Context) error {
	audioData, _ := c.FormFile("audio")

	isAudio := utils.CheckAudioFile(audioData)
	if !isAudio {
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
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"public_song_lyrics": resData,
	})
}
