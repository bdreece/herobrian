package minecraft

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"slices"

	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type (
	HostInfo struct {
		ID           string        `json:"-"`
		Architecture string        `json:"arch"`
		DNSName      string        `json:"-"`
		Image        *ImageInfo    `json:"image"`
		IPAddress    string        `json:"-"`
		Platform     string        `json:"platform"`
		Processor    ProcessorInfo `json:"processor"`
	}

	ImageInfo struct {
		Type string `json:"type"`
		// Memory size in MiB.
		Memory *int64 `json:"memory"`
		// Storage space in GiB.
		Storage *int64 `json:"storage"`
		Network string `json:"network"`
	}

	ProcessorInfo struct {
		CoreCount      *int32 `json:"cores"`
		ThreadsPerCore *int32 `json:"threads"`
	}

	HostProvider interface {
		Hosts(ctx context.Context, ids ...string) ([]*HostInfo, error)
		StartHosts(ctx context.Context, ids ...string) error
		StopHosts(ctx context.Context, ids ...string) error
		RestartHosts(ctx context.Context, ids ...string) error
	}

	HostControllerParams struct {
		fx.In

		Provider HostProvider
		Logger   *slog.Logger
	}

	HostController struct {
		provider HostProvider
		logger   *slog.Logger
		hosts    map[string]HostConfig
	}
)

var ErrHostNotFound = errors.New("minecraft: host not found")

func NewHostController(p HostControllerParams) (*HostController, error) {
	hosts := map[string]HostConfig{}
	if err := viper.UnmarshalKey("minecraft:hosts", &hosts); err != nil {
		return nil, err
	}

	slog.Debug("unmarshaled hosts", "hosts", hosts)

	controller := HostController{
		p.Provider,
		p.Logger.With("scope", "minecraft.HostController"),
		hosts,
	}

	return &controller, nil
}

func (h *HostController) Routes() []echo.Route {
	return []echo.Route{
		{Method: http.MethodGet, Path: "/host", Handler: h.List},
	}
}

func (h *HostController) List(c *echo.Context) error {
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

	result := maps.Collect(func(yield func(string, *HostInfo) bool) {
		for name := range h.hosts {
			hostIndex := slices.IndexFunc(hosts, func(host *HostInfo) bool {
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
