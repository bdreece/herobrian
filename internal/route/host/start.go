package host

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Controller) Start(c *echo.Context) error {
	client, _ := FromContext(c)

	if err := client.Start(c.Request().Context()); err != nil {
		return echo.ErrBadGateway.Wrap(err)
	}

	return c.NoContent(http.StatusAccepted)
}
