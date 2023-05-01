package controllers

import (
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/Lirikku/models"
	"github.com/Lirikku/utils"
	"github.com/labstack/echo/v4"
)
func SearchTermSongLyrics(c echo.Context) error {
	term := c.QueryParam("term")
	offset := c.QueryParam("offset")
	
	offsetInt, err := utils.CheckOffset(offset)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "offset must be a number",
		})
	}

	res, err := utils.RequestShazamSearchTerm(term, strconv.Itoa(offsetInt), "artists,songs", "5")
	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "failed to get data",
		})
	}

	keys := res.GetKeys()
	
	var resPublicSongLyrics = make([]models.ResponsePublicSongLyric, len(keys))

	var wg sync.WaitGroup
	for i, key := range keys {
		wg.Add(1)
		go func (i int, key string) {
			defer wg.Done()

			res, err := utils.RequestShazamSearchKey(key)
			
			if err != nil {
				return
			}
			
			resPublicSongLyrics[i] = res.GetInResponsePublicSongLyric()
			
		}(i, key)

	}
	wg.Wait()

	next := utils.GenerateNextLink(c, len(resPublicSongLyrics), url.Values{
		"term": {term},
		"offset": {strconv.Itoa(offsetInt + 5)},
	}.Encode())

	return c.JSON(http.StatusOK, echo.Map{
		"next": next,
		"public_song_lyrics": resPublicSongLyrics,
	})

}

func SearchAudioSongLyrics(c echo.Context) error {
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

	key, err := utils.RequestShazamSearchAudio(rawBases64)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "failed to get data",
		})
	}

	if key == "" {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "song not found",
		})
	}

	res, err := utils.RequestShazamSearchKey(key)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "failed to get data",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"public_song_lyrics": res.GetInResponsePublicSongLyric(),
	})
}
