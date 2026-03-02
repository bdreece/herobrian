package route

import "github.com/labstack/echo/v5"

type Router interface {
	Routes() []echo.Route
}
