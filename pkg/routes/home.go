package routes

import (
	"log/slog"
	"net/http"

	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/labstack/echo/v4"
)

func Home(client *linode.Client, log *slog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
        log.Info("fetching instance status...")
        status, err := client.InstanceStatus(c.Request().Context())
        if err != nil {
            return echo.NewHTTPError(http.StatusFailedDependency, err.Error())
        }

        log.Info("got instance status", slog.String("status", status))

		return c.Render(http.StatusOK, "home.gotmpl", echo.Map{
            "URL": "minecraft.bdreece.dev",
            "Status": status,
        })
	}
}
