package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GenerateNextLink(c echo.Context, offset, lenSongLyrics int) string {
	var next string
	if lenSongLyrics < 5 {
		next = ""
	} else {
		next = "http://"+ c.Request().Host + c.Request().URL.Path + "?offset=" + strconv.Itoa(offset + 5)
	}
	
	return next
}