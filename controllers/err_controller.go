package controllers

import (
	"net/http"

	"github.com/Lirikku/utils"
	"github.com/Lirikku/view"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
	}
	if err = utils.Render(c, code, view.Error(he)); err != nil {
		c.Logger().Error(err)
	}	
}
