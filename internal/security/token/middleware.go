package token

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

const (
	claimsContextKey string = "claims"
	tokenContextKey  string = "token"
)

func Authenticate(decoder Decoder[*AccessClaims]) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			header := c.Request().Header.Get("authorization")
			jwt := strings.TrimPrefix(header, "Bearer ")

			claims, token, err := decoder.Decode(jwt)
			if err == nil {
				c.Set(claimsContextKey, claims)
				c.Set(tokenContextKey, token)
			}

			return next(c)
		}
	}
}

func Verify() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token, ok := c.Get(tokenContextKey).(*jwt.Token)
			if !ok || !token.Valid {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}
