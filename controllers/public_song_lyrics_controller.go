package controllers

import (
	"fmt"
	"net/http"

	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	"github.com/Lirikku/utils"
	"github.com/Lirikku/view"
	view_component "github.com/Lirikku/view/component"
	"github.com/labstack/echo/v4"
)

type PublicSongLyrics struct {
	service   services.IPublicSongLyricsService
	myService services.IMySongLyricsService
}

func NewPublicSongLyricsController(service services.IPublicSongLyricsService, myService services.IMySongLyricsService) *PublicSongLyrics {
	return &PublicSongLyrics{service, myService}
}

func (pub *PublicSongLyrics) SongLyricsView(c echo.Context) error {	
	return utils.Render(c, http.StatusOK, view.PubSongs(c))
}

func (pub *PublicSongLyrics) SearchSongsByTerm(c echo.Context) error {
	term := c.QueryParam("term")

	currentPage, pageSize, offset := utils.GetPageSizeAndOffset(c)
	res, keyPubs, err := pub.service.SearchSongsByTermShazam(term, "artists,songs", offset, pageSize)
	if err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-results")
		c.Response().Header().Set("HX-Swap", "innerHTML")
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	if len(c.Request().Header.Values("HX-Request")) > 0 && c.Request().Header.Values("HX-Request")[0] == "true" {
		logged, ok := c.Get("auth").(bool)
		if ok && logged {
			user := c.Get("user").(models.UserJWTDecode)
			if err := pub.myService.CheckKeyPub(user.ID, res, keyPubs); err != nil {
				return utils.Render(c, http.StatusOK, view.Error(echo.ErrInternalServerError))
			}
		}
		c.Set("current_page", currentPage)
		c.Set("next_name", c.Path())
		c.Set("term", term)
		c.Set("res", res)
		return utils.Render(c, http.StatusOK, view.PubResultSongs(c))
	}

	return c.NoContent(http.StatusNoContent)
}

func (pub *PublicSongLyrics) SaveSong(c echo.Context) error {
	req := new(models.PublicSaveSong)
	c.Bind(req)
	if req.ArtistNames == "" || req.Title == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "empty field not allowed",
		})
	}

	detailedSong, err := pub.service.GetSongDetail(req.ArtistNames, req.Title)
	if err != nil {
		fmt.Println("save pub song -> ",err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}
	user := c.Get("user").(models.UserJWTDecode)
	if err := pub.myService.SaveSongLyric(user.ID, models.SongLyricWrite{
		ArtistNames: detailedSong.ArtistNames,
		Title:       detailedSong.Title,
		Key:         req.Key,
		Lyric:       detailedSong.Lyric,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "internal server error",
		})
	}

	if len(c.Request().Header.Values("HX-Request")) > 0 && c.Request().Header.Values("HX-Request")[0] == "true" {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "successfully",
		})
	}

	return c.NoContent(http.StatusNoContent)

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
		return utils.Render(c, http.StatusOK, view.PubSongDetail(res))
	}

	return c.NoContent(http.StatusNoContent)
}

func (pub *PublicSongLyrics) SearchAudioSongLyric(c echo.Context) error {
	audioData, _ := c.FormFile("audio")

	q := new(models.ReqSearchAudio)

	c.Bind(q)

	fmt.Println("query -> ", q)

	if audioData == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "data is empty")
	}

	isAudio := utils.CheckAudioFile(audioData)

	if !isAudio {
		return c.String(http.StatusBadRequest, "invalid file type. please upload an audio file")
	}

	if audioData.Size > 500000 {
		return c.String(http.StatusBadRequest, "audio size must be less than 500kb")
	}

	rawBases64 := utils.Audio2RawBase64(audioData)

	resData, err := pub.service.SearchSongLyricByAudioRapidShazam(rawBases64, q)
	// resData = models.RapidShazamSearchAudioResponse{}

	if err != nil {
		return utils.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	hreq := c.Request().Header.Values("HX-Request")
	if len(hreq) > 0 && hreq[0] == "true" {
		if resData.Track.Title == "" {
			return c.String(http.StatusBadRequest, "couldn't quite catch that")
		}

		detail, _ := pub.service.GetSongDetail(resData.Track.Subtitle, resData.Track.Title)
		return utils.Render(c, http.StatusPartialContent, view_component.SongsDetail(utils.Convert2Map(detail)))
	}

	return c.NoContent(http.StatusNoContent)
}
