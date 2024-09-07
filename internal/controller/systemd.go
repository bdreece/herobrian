package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/bdreece/herobrian/pkg/cron"
	"github.com/bdreece/herobrian/pkg/systemd"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/singleflight"
)

var (
	ErrInstanceNotFound = errors.New("systemd unit instance not found")
)

type Systemd struct {
	services *systemd.ServiceFactory
	logger   *slog.Logger
	group    singleflight.Group
}

type systemdModel struct {
	Instance string `param:"instance" validate:"required"`
}

func (controller *Systemd) Enable(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err := svc.Enable(c.Request().Context()); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`
        <span
            id="%s-status"
            class="rounded-full bg-neutral-200 p-2"
            sse-swap="status"
            hx-swap="outerHTML"
        >
            enabling...
        </span>
    `, svc.Unit().Instance))
}

func (controller *Systemd) Disable(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err := svc.Disable(c.Request().Context()); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`
        <span
            id="%s-status"
            class="rounded-full bg-neutral-200 p-2"
            sse-swap="status"
            hx-swap="outerHTML"
        >
            disabling...
        </span>
    `, svc.Unit().Instance))
}

func (controller *Systemd) Start(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err := svc.Start(c.Request().Context()); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`
        <span
            id="%s-status"
            class="rounded-full bg-neutral-200 p-2"
            sse-swap="status"
            hx-swap="outerHTML"
        >
            starting...
        </span>
    `, svc.Unit().Instance))
}

func (controller *Systemd) Stop(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err := svc.Stop(c.Request().Context()); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`
        <span
            id="%s-status"
            class="rounded-full bg-neutral-200 p-2"
            sse-swap="status"
            hx-swap="outerHTML"
        >
            stopping...
        </span>
    `, svc.Unit().Instance))
}

func (controller *Systemd) Restart(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err = svc.Restart(c.Request().Context()); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`
        <span
            id="%s-status"
            class="rounded-full bg-neutral-200 p-2"
            sse-swap="status"
            hx-swap="outerHTML"
        >
            restarting...
        </span>
    `, svc.Unit().Instance))
}

func (controller *Systemd) SSE(c echo.Context) error {
	var (
		buf      bytes.Buffer
		interval = 15 * time.Second
	)

	model := new(systemdModel)
	if err := c.Bind(model); err != nil {
		return err
	}

	if err := c.Validate(model); err != nil {
		return err
	}

	service, err := controller.services.Create(model.Instance)
	if err != nil {
		return err
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for range cron.Tick(c.Request().Context(), time.Now().Round(interval), interval) {
        key := fmt.Sprintf("systemd-status-%s", model.Instance)
		result := <-controller.group.DoChan(key, func() (interface{}, error) {
			controller.logger.Info("requesting instance state", slog.String("instance", model.Instance))
			return service.Status(c.Request().Context())
		})

		if result.Err != nil {
            controller.logger.Error("failed to get instance status",
                slog.String("instance", model.Instance))

			return result.Err
		}

		if result.Shared {
			controller.logger.Info("instance state result was shared!")
		}

		status := *result.Val.(*systemd.Status)
        controller.logger.Info("got instance status",
            slog.String("instance", model.Instance),
            slog.String("status", status.String()))

		_, err = systemdStatusEvent(model.Instance, status).WriteTo(&buf)
		if err != nil {
			return err
		}

		if _, err = io.Copy(w, &buf); err != nil {
			return err
		}

		w.Flush()
		buf.Reset()
	}

	return nil
}

func (controller *Systemd) resolveService(c echo.Context) (*systemd.Service, error) {
	model := new(systemdModel)
	if err := c.Bind(model); err != nil {
		return nil, err
	}

	if err := c.Validate(model); err != nil {
		return nil, err
	}

	svc, err := controller.services.Create(model.Instance)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, err)
	}

	return svc, nil
}

func NewSystemd(services *systemd.ServiceFactory, logger *slog.Logger) *Systemd {
	return &Systemd{
		services: services,
		logger:   logger,
	}
}

func systemdStatusEvent(instance string, status systemd.Status) event {
	var variant string
	switch status {
	case systemd.StatusActiveRunning, systemd.StatusEnabled:
		variant = "bg-secondary"
	case systemd.StatusFailed, systemd.StatusActiveExited, systemd.StatusDisabled, systemd.StatusInactive:
		variant = "bg-red-300"
	case systemd.StatusActiveWaiting:
		variant = "bg-primary"
	default:
		variant = "bg-neutral-200"
	}

	return event{
		Event: "status",
		Data: fmt.Sprintf(`
            <span
                id="%s-status"
                class="rounded-full %s p-2"
                sse-swap="status"
                hx-swap="outerHTML"
            >
                %s
            </span>
        `, instance, variant, status),
	}
}
