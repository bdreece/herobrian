package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewMiddleware(opts *Options) echo.MiddlewareFunc {
    return middleware.BasicAuth(func(user, pass string, ctx echo.Context) (bool, error) {
        return user == opts.Username && pass == opts.Password, nil
    })
}
