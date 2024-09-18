package utils

import (
	"github.com/labstack/echo/v4"
)

func GetContext[T any](c echo.Context, key string) T {
	return c.Get(key).(T)
}

func SetContext[T any](c echo.Context, key string, val T) {
	c.Set(key, val)
}

func ErrResponse(c echo.Context, status int, msg string, setOthers ...func(c echo.Context)) error {
	c.Response().Header().Set("HX-Swap", "innerHTML")
	c.Response().Header().Set("HX-Retarget", "#error-results")
	for _, v := range setOthers {
		v(c)
	}
	return c.String(status, msg)
}
