package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Lirikku/models"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JWTClaims)
		},
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			return jwt.ParseWithClaims(auth, new(models.JWTClaims), func(token *jwt.Token) (interface{}, error) {
				if token.Claims.(*models.JWTClaims).ExpiresAt.Time.Before(time.Now()) {
					return nil, fmt.Errorf("token expired")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
		},
		ErrorHandler: func(c echo.Context, err error) error {
			if err.Error() != "token expired" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		},
	})
}