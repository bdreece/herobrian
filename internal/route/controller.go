package route

import "github.com/labstack/echo/v5"

type Controller interface {
	Routes() []echo.Route
}
