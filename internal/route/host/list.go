package host

import (
	"net/http"
	"slices"

	"github.com/labstack/echo/v5"
)

func (h *Controller) List(c *echo.Context) error {
	ids := slices.Collect(func(yield func(string) bool) {
		for key := range h.hosts {
			if !yield(h.hosts[key].ID) {
				return
			}
		}
	})

	hosts, err := h.describer.DescribeHosts(c.Request().Context(), ids...)
	if err != nil {
		return echo.ErrBadGateway.Wrap(err)
	}

	c.Response().Header().Add("Cache-Control", "public, max-age=1800")

	return c.JSON(http.StatusOK, hosts)
}
