package utils

import (
	"github.com/labstack/echo/v4"
)

func GenerateNextLink(c echo.Context, lenData int, param string) string {
	var next string
	if lenData < 5 {
		next = ""
	} else {
		next = "http://"+ c.Request().Host + c.Request().URL.Path + "?" + param
	}
	
	return next
}