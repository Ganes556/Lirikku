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

	shazamMetadata, err := utils.RequestShazamMetadata(term, strconv.Itoa(offsetInt), "artists,songs", "5")
	
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "failed to get data",
		})
	}

	keys := shazamMetadata.GetKeys()
	titles := shazamMetadata.GetTitles()
	artists := shazamMetadata.GetArtists()
	
	var resPublicSongLyrics = make([]models.ResponsePublicSongLyric, len(keys))

	var wg sync.WaitGroup
	for i, key := range keys {
		wg.Add(1)
		go func (i int, key string) {
			defer wg.Done()

			shazamLyric, err := utils.RequestShazamLyric(key)
			
			if err != nil {
				return
			}

			if i < len(titles) {
				resPublicSongLyrics[i].Title = titles[i]
			}

			if i < len(artists) {
				resPublicSongLyrics[i].ArtistNames = artists[i]
			}
			
			lyric := shazamLyric.GetLyrics()
			
			resPublicSongLyrics[i].Lyric = lyric
			
			
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

func SearchAudioSongLyric(c echo.Context) error {
	return nil
}
