package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bdreece/herobrian/pkg/systemd"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var (
	ErrInstanceNotFound = errors.New("systemd unit instance not found")
)

type SystemdParams struct {
	fx.In

	Emitter  systemd.Emitter
	Services *systemd.ServiceFactory
}

type Systemd struct {
	emitter  systemd.Emitter
	services *systemd.ServiceFactory
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

	controller.emitter.Publish(svc.Unit().Instance, systemd.StatusEnabled)
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

	controller.emitter.Publish(svc.Unit().Instance, systemd.StatusDisabled)
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
	model := new(systemdModel)
	if err := c.Bind(model); err != nil {
		return err
	}

	if err := c.Validate(model); err != nil {
		return err
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	statusch := make(chan systemd.Status, 1)
	sub := controller.emitter.Subscribe(model.Instance, statusch)
	defer controller.emitter.Unsubscribe(model.Instance, sub)

	var buf bytes.Buffer
	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case status := <-statusch:
			_, err := systemdStatusEvent(model.Instance, status).WriteTo(&buf)
			if err != nil {
				return err
			}

			if _, err = io.Copy(w, &buf); err != nil {
				return err
			}

			w.Flush()
			buf.Reset()
		}
	}
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

func NewSystemd(p SystemdParams) *Systemd {
	return &Systemd{
		services: p.Services,
		emitter:  p.Emitter,
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
