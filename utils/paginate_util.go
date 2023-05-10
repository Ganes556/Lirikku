package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetPageSizeAndOffset(c echo.Context) (int, int) {
	
	page, _ := strconv.Atoi(c.QueryParam("page"))

	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	
	switch {
    case pageSize > 15:
      pageSize = 15
    case pageSize <= 0:
      pageSize = 5
  }

	return pageSize, (page - 1) * pageSize
}