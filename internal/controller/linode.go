package controller

import (
	"bytes"
	"fmt"
	"io"
    "log/slog"
	"net/http"
	"time"

	"github.com/bdreece/herobrian/pkg/cron"
	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/singleflight"
)

var (
	linodeRunningEvent = event{
		Event: "running",
		Data: `
            <button
                class="btn btn-primary"
                hx-post="/linode/shutdown"
                hx-swap="outerHTML"
            >
                Shutdown
            </button> 
        `,
	}

	linodeOfflineEvent = event{
		Event: "offline",
		Data: `
            <button
                class="btn btn-secondary"
                hx-post="/linode/boot"
                hx-swap="outerHTML"
            >
                Boot
            </button>
        `,
	}

	spinner = `
        <div
            class="sk-cube-grid"
            sse-swap="running,offline"
            hx-swap="outerHTML"
        >
            <div class="sk-cube sk-cube1"></div>
            <div class="sk-cube sk-cube2"></div>
            <div class="sk-cube sk-cube3"></div>
            <div class="sk-cube sk-cube4"></div>
            <div class="sk-cube sk-cube5"></div>
            <div class="sk-cube sk-cube6"></div>
            <div class="sk-cube sk-cube7"></div>
            <div class="sk-cube sk-cube8"></div>
            <div class="sk-cube sk-cube9"></div>
        </div>
    `
)

const key string = "linode"

type Linode struct {
	client linode.Client
    logger *slog.Logger
	group  singleflight.Group
}

func (controller *Linode) Boot(c echo.Context) error {
	err := controller.client.BootInstance(c.Request().Context())
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, spinner)
}

func (controller *Linode) Reboot(c echo.Context) error {
	err := controller.client.RebootInstance(c.Request().Context())
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, spinner)
}

func (controller *Linode) Shutdown(c echo.Context) error {
	err := controller.client.ShutdownInstance(c.Request().Context())
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, spinner)
}

func (controller *Linode) SSE(c echo.Context) error {
	var (
		buf      bytes.Buffer
		interval = 2 * time.Second
	)

	c.Logger().Info("subscribing to server-sent events")

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Flush()

	for range cron.Tick(c.Request().Context(), time.Now().Round(interval), interval) {
        result := <-controller.group.DoChan(key, func() (interface{}, error) {
            controller.logger.Info("requesting linode status...")
			return controller.client.InstanceStatus(c.Request().Context())
		})

		if result.Err != nil {
			return result.Err
		}

        if result.Shared {
            controller.logger.Info("linode state result was shared!")
        }

		status := *result.Val.(*linode.Status)
        controller.logger.Debug("got status", slog.String("status", status.String()))
		_, _ = linodeStatusEvent(status).WriteTo(&buf)

		if status == linode.StatusRunning {
			_, _ = linodeRunningEvent.WriteTo(&buf)
		} else if status == linode.StatusOffline {
			_, _ = linodeOfflineEvent.WriteTo(&buf)
		}

		if _, err := io.Copy(w, &buf); err != nil {
			return err
		}

		w.Flush()
		buf.Reset()
	}

    return nil
}

func NewLinode(client linode.Client, logger *slog.Logger) *Linode {
	return &Linode{
        client: client,
        logger: logger,
    }
}

func linodeStatusEvent(status linode.Status) event {
	var title string
	switch status {
	case linode.StatusRunning:
		title = "The linode instance is running"
	case linode.StatusOffline:
		title = "The linode instance is offline"
	case linode.StatusRebooting:
		title = "The linode instance is currently rebooting"
	case linode.StatusShuttingDown:
		title = "The linode instance is currently shutting down"
	case linode.StatusBooting:
		title = "The linode instance is currently booting"
	}

	return event{
		Event: "status",
		Data: fmt.Sprintf(`
            <span title="%s">%s</span>
        `, title, status),
	}
}
