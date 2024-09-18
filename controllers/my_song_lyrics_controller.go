package controllers

import (
	"fmt"
	"net/http"

	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	"github.com/Lirikku/utils"
	"github.com/Lirikku/view"
	"github.com/labstack/echo/v4"
)

type MySongLyrics struct {
	service services.IMySongLyricsService
}

func NewMySongLyricsController(service services.IMySongLyricsService) *MySongLyrics {
	return &MySongLyrics{service}
}

func (my *MySongLyrics) GetSongLyrics(c echo.Context) error {

	user, _ := c.Get("user").(models.UserJWTDecode)
	// if !ok {
	// }
	term := c.QueryParam("term")
	currentPage, pageSize, offset := utils.GetPageSizeAndOffset(c)

	res, err := my.service.GetSongLyrics(user.ID, offset, pageSize)

	if err != nil {
		return utils.Render(c, http.StatusInternalServerError, view.Error(echo.ErrInternalServerError))
	}

	c.Set("current_page", currentPage)
	c.Set("next_name", c.Echo().URL(my.GetSongLyric))
	c.Set("term", term)
	c.Set("res", res)

	if len(c.Request().Header.Values("HX-Request")) > 0 && c.Request().Header.Values("HX-Request")[0] == "true" {
		if currentPage != 1 {
			return utils.Render(c, http.StatusOK, view.MyResultSongs(c))
		}
	}

	return utils.Render(c, http.StatusOK, view.MySongs(c))
	// return c.JSON(http.StatusOK, echo.Map{
	// 	"my_song_lyrics": resSongLyrics,
	// })
}

func (my *MySongLyrics) GetSongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)

	idSongLyric := c.Param("id")

	idSongLyricInt := utils.CheckId(idSongLyric)

	if idSongLyricInt == -1 {
		// return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
		// 	"message": "id must be a number and greater than 0",
		// })
		return utils.ErrResponse(c, http.StatusBadRequest, "id must be a number and greater than 0")
	}

	res, err := my.service.GetSongLyric(idSongLyricInt, user.ID)

	if err != nil {
		return utils.ErrResponse(c, http.StatusNotFound, "song lyric not found")
		// return echo.NewHTTPError(http.StatusNotFound, echo.Map{
		// 	"message": "song lyric not found",
		// })
	}

	if len(c.Request().Header.Values("HX-Request")) > 0 && c.Request().Header.Values("HX-Request")[0] == "true" {
		qPartial := c.QueryParam("partial")
		if qPartial == "dialog-edit" {
			c.Response().Header().Set("HX-Swap", "innerHTML")
			c.Response().Header().Set("HX-Retarget", "#form-edit-content")
			return utils.Render(c, http.StatusOK, view.DialongEditInput(res))
		}
		return utils.Render(c, http.StatusOK, view.MySongDetail(res))
	}

	// return c.JSON(http.StatusOK, echo.Map{
	// 	"my_song_lyrics": resSongLyric,
	// })
	return c.NoContent(http.StatusNoContent)
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

	currentPage, pageSize, offset := utils.GetPageSizeAndOffset(c)

	term := c.QueryParam("term")

	res, err := my.service.SearchSongLyrics(user.ID, term, offset, pageSize)

	if err != nil {
		return utils.ErrResponse(c, http.StatusInternalServerError, "internal server error")
	}

	if len(res) == 0 {
		return utils.ErrResponse(c, http.StatusNotFound, "song not found", func(c echo.Context) {
			c.Response().Header().Set("HX-Swap", "innerHTML")
			c.Response().Header().Set("HX-Retarget", "#search-results")
		})
	}

	if len(c.Request().Header.Values("HX-Request")) > 0 && c.Request().Header.Values("HX-Request")[0] == "true" {
		c.Set("current_page", currentPage)
		c.Set("next_name", c.Path())
		c.Set("term", term)
		c.Set("res", res)
		return utils.Render(c, http.StatusOK, view.MyResultSongs(c))
	}

	return c.NoContent(http.StatusNoContent)
}

func (my *MySongLyrics) DeleteSongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)

	idSongLyric := c.Param("id")

	idSongLyricInt := utils.CheckId(idSongLyric)

	if idSongLyricInt <= 0 {
		fmt.Println(idSongLyricInt)
		return utils.ErrResponse(c, http.StatusInternalServerError, "internal server error")
	}

	_, err := my.service.GetSongLyric(idSongLyricInt, user.ID)

	if err != nil {
		return utils.ErrResponse(c, http.StatusNotFound, "song not found", func(c echo.Context) {
			c.Response().Header().Set("HX-Swap", "innerHTML")
			c.Response().Header().Set("HX-Retarget", "#search-results")
		})
	}

	err = my.service.DeleteSongLyric(idSongLyricInt, user.ID)

	if err != nil {
		return utils.ErrResponse(c, http.StatusInternalServerError, "internal server error")
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (my *MySongLyrics) UpdateSongLyric(c echo.Context) error {

	user := c.Get("user").(models.UserJWTDecode)

	var req models.SongLyricWrite

	c.Bind(&req)

	idSongLyricInt := utils.CheckId(req.ID)

	if idSongLyricInt <= 0 {
		return utils.ErrResponse(c, http.StatusInternalServerError, "internal server error")
	}

	_, err := my.service.GetSongLyric(idSongLyricInt, user.ID)

	if err != nil {
		return utils.ErrResponse(c, http.StatusNotFound, "song not found", func(c echo.Context) {
			c.Response().Header().Set("HX-Retarget", "#search-results")
		})
	}

	err = my.service.UpdateSongLyric(idSongLyricInt, user.ID, req)

	if err != nil {
		return utils.ErrResponse(c, http.StatusInternalServerError, "internal server error")
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}
