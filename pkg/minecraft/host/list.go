package host

import (
	"maps"
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

	hosts, err := h.provider.Hosts(c.Request().Context(), ids...)
	if err != nil {
		return echo.ErrBadGateway.Wrap(err)
	}

	result := maps.Collect(func(yield func(string, *Info) bool) {
		for name := range h.hosts {
			hostIndex := slices.IndexFunc(hosts, func(host *Info) bool {
				return host.ID == h.hosts[name].ID
			})

			if hostIndex < 0 {
				continue
			}

			if !yield(name, hosts[hostIndex]) {
				return
			}
		}
	})

	c.Response().Header().Add("Cache-Control", "public, max-age=1800")

	return c.JSON(http.StatusOK, result)
}
