package controller

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/bdreece/herobrian/pkg/systemd"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

const (
	systemdEnabledEvent  = "event: enabled\ndata:\n\n"
	systemdDisabledEvent = "event: disabled\ndata:\n\n"
	systemdStartedEvent  = "event: started\ndata:\n\n"
	systemdStoppedEvent  = "event: stopped\ndata:\n\n"
)

var (
	ErrInstanceNotFound = errors.New("systemd unit instance not found")
)

type SystemdParams struct {
	fx.In

	Emitter  systemd.Emitter
	Services *systemd.ServiceFactory
	Logger   *slog.Logger
}

type Systemd struct {
	emitter  systemd.Emitter
	services *systemd.ServiceFactory
	logger   *slog.Logger
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
	return c.NoContent(http.StatusOK)
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
	return c.NoContent(http.StatusOK)
}

func (controller *Systemd) Start(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err := svc.Start(c.Request().Context()); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (controller *Systemd) Stop(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err := svc.Stop(c.Request().Context()); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (controller *Systemd) Restart(c echo.Context) error {
	svc, err := controller.resolveService(c)
	if err != nil {
		return err
	}

	if err = svc.Restart(c.Request().Context()); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
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

	ch := make(chan systemd.Status, 1)
	sub := controller.emitter.Subscribe(model.Instance, ch)
	defer controller.emitter.Unsubscribe(model.Instance, sub)

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case status := <-ch:
			_, err := io.WriteString(w, systemdStatusEvent(status))
			if err != nil {
				return err
			}

			w.Flush()
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
		logger:   p.Logger,
	}
}

func systemdStatusEvent(status systemd.Status) string {
	return fmt.Sprintf("event: status\ndata: %s\n\n", status)
}
