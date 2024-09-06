package middleware

import (
	"github.com/bdreece/herobrian/pkg/identity"
	"github.com/labstack/echo/v4"
)

const ClaimsContextKey string = "claims"

func Authenticate(authenticator identity.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := authenticator.Authenticate(c)
			if err != nil {
				c.Logger().Warnf("failed to authenticate user: %v", err)
			}

			c.Set(ClaimsContextKey, claims)
			return next(c)
		}
	}
}
