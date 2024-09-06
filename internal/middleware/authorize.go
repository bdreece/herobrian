package middleware

import (
	"net/http"

	"github.com/bdreece/herobrian/pkg/identity"
	"github.com/labstack/echo/v4"
)

func Authorize(authorizers ...identity.Authorizer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(ClaimsContextKey).(*identity.ClaimSet)
			if !ok || claims == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "failed to get claims from request context")
			}

			for _, authorizer := range authorizers {
				if err := authorizer.Authorize(claims); err != nil {
					return echo.NewHTTPError(http.StatusForbidden, err.Error())
				}
			}

			return next(c)
		}
	}
}
