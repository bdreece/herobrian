package host

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Controller) Stop(c *echo.Context) error {
	hostParam := c.Param("host")
	host, ok := h.hosts[hostParam]
	if !ok {
		return echo.ErrBadRequest
	}

	if err := h.provider.StopHosts(c.Request().Context(), host.ID); err != nil {
		return echo.ErrBadGateway.Wrap(err)
	}

	return c.NoContent(http.StatusAccepted)
}
