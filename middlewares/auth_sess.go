package middlewares

import (
	"github.com/Lirikku/configs"
	"github.com/labstack/echo/v4"
)

func Authorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			s, err := configs.Store.Get(c.Request(), "session")
			if err != nil {
				c.Response().Header().Set("HX-Redirect", "/auth/login")
				// return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			if auth, ok := s.Values["auth"].(bool); !ok || !auth {
				c.Response().Header().Set("HX-Redirect", "/auth/login")
				// return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			return next(c)
		}
	}
}
