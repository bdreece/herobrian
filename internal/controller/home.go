package controller

import (
	"log/slog"
	"net/http"

	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/bdreece/herobrian/pkg/systemd"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type (
	Home struct {
		client   linode.Client
		services *systemd.ServiceFactory
		logger   *slog.Logger
	}

	HomeParams struct {
		fx.In

		LinodeClient   linode.Client
		ServiceFactory *systemd.ServiceFactory
		Logger         *slog.Logger
	}
)

func (controller *Home) RenderIndex(c echo.Context) error {
	c.Logger().Info("fetching instance status...")
	status, err := controller.client.InstanceStatus(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusFailedDependency, err.Error())
	}

	c.Logger().Info("got instance status", slog.String("status", status.String()))

	return c.Render(http.StatusOK, "home.gotmpl", echo.Map{
		"URL":    "minecraft.bdreece.dev",
		"Units":  controller.services.Units(),
		"Status": status.String(),
	})
}

func NewHome(p HomeParams) *Home {
	return &Home{
		client:   p.LinodeClient,
		services: p.ServiceFactory,
		logger:   p.Logger,
	}
}
