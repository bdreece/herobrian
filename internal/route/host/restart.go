package host

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Controller) Restart(c *echo.Context) error {
	client, _ := FromContext(c)

	if err := client.Restart(c.Request().Context()); err != nil {
		return echo.ErrBadGateway.Wrap(err)
	}

	return c.NoContent(http.StatusAccepted)
}
