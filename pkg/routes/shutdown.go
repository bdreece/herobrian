package routes

import (
	"net/http"

	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/labstack/echo/v4"
)

func Shutdown(client *linode.Client) echo.HandlerFunc {
    return func(c echo.Context) error {
        if err := client.ShutdownInstance(c.Request().Context()); err != nil {
            return echo.NewHTTPError(http.StatusFailedDependency, err.Error())
        }

        return c.HTML(http.StatusOK, `
            <span
                class="rounded-full bg-slate-200 text-center"
                sse-swap="running,offline"
                hx-swap="outerHTML"
            >
                Shutting down...
            </span>
        `)
    }
}
