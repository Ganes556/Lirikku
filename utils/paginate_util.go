package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetPageSizeAndOffset(c echo.Context) (int, int, int) {

	current, _ := strconv.Atoi(c.QueryParam("page"))

	if current <= 0 {
		current = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	switch {
	case pageSize > 20:
		pageSize = 20
	case pageSize <= 0:
		pageSize = 10
	}

	return current, pageSize, (current - 1) * pageSize
}
