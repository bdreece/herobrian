package identity

import "github.com/labstack/echo/v4"

func NewActivateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var model struct {
			DisplayName string `form:"displayName"`
		}
	}
}
