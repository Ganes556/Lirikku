package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
	"github.com/Lirikku/utils"
	"github.com/Lirikku/view"
	"github.com/labstack/echo/v4"
)

func Authorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			s, _ := configs.Store.Get(c.Request(), "session")
			path := c.Request().URL.Path
			auth, ok := s.Values["auth"].(bool)
			switch path {
			case "/my", "/save":
				if !ok || !auth {
					c.Redirect(http.StatusPermanentRedirect, "/auth/login")
					return next(c)
				}
			case "/auth/login", "/auth/register":
				if ok && auth {
					c.Redirect(http.StatusPermanentRedirect, "/")
				}
			default:
				if !ok {
					c.Set("auth", false)
				} else {
					c.Set("auth", auth)
				}
			}
			if c.Get("user") == nil || c.Get("auth") == false {
				jwtUserData, ok := s.Values["user"].(string)
				if ok {
					var user models.UserJWTDecode
					if err := json.Unmarshal([]byte(jwtUserData), &user); err != nil {
						return utils.Render(c, http.StatusInternalServerError, view.Error(echo.ErrInternalServerError))
					}
					c.Set("user", user)
				}
			}
			return next(c)
		}
	}
}
