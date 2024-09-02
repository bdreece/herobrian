package controller

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
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

type LinodeParams struct {
	fx.In

	Client  linode.Client
	Emitter linode.Emitter
	Logger  *slog.Logger
}

type Linode struct {
	client  linode.Client
	emitter linode.Emitter
	logger  *slog.Logger
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
	controller.logger.Info("subscribing to server-sent events")
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Flush()

	ch := make(chan linode.Status, 1)
	sub := controller.emitter.Subscribe(linode.TopicStatus, ch)
	defer controller.emitter.Unsubscribe(linode.TopicStatus, sub)

	var buf bytes.Buffer

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case status := <-ch:
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
	}
}

func NewLinode(p LinodeParams) *Linode {
	return &Linode{
		client:  p.Client,
		emitter: p.Emitter,
		logger:  p.Logger,
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
