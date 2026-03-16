package host

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/r3labs/sse/v2"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type (
	ControllerParams struct {
		fx.In

		Provider Provider
		SSE      *sse.Server
		Logger   *slog.Logger
	}

	Controller struct {
		provider Provider
		sse      *sse.Server
		logger   *slog.Logger
		hosts    map[string]Config
	}
)

var ErrHostNotFound = errors.New("minecraft: host not found")

func NewController(p ControllerParams) (*Controller, error) {
	hosts := map[string]Config{}
	if err := viper.UnmarshalKey("minecraft:hosts", &hosts); err != nil {
		return nil, err
	}

	slog.Debug("unmarshaled hosts", "hosts", hosts)

	controller := Controller{
		p.Provider,
		p.SSE,
		p.Logger.With("scope", "minecraft.HostController"),
		hosts,
	}

	return &controller, nil
}

func (h *Controller) Routes() []echo.Route {
	return []echo.Route{
		{Method: http.MethodGet, Path: "/host", Handler: h.List},
		{Method: http.MethodPost, Path: "/host/:host/start", Handler: h.Start},
		{Method: http.MethodPost, Path: "/host/:host/stop", Handler: h.Stop},
		{Method: http.MethodPost, Path: "/host/:host/restart", Handler: h.Restart},
	}
}
