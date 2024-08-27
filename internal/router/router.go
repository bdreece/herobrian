package router

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/internal/controller"
)

type Params struct {
	fx.In

	Renderer  echo.Renderer
	Validator echo.Validator
	Logger    *slog.Logger
	Options   *Options
}

type Router struct {
	*echo.Echo
}

func New(p Params) Router {
	e := echo.New()
	e.Renderer = p.Renderer
	e.Validator = p.Validator

	e.Use(
		middleware.BodyLimit("4M"),
		middleware.Decompress(),
		middleware.Gzip(),
		middleware.CSRF(),
		middleware.Secure(),
		middleware.Static(p.Options.StaticDirectory),
		middleware.Static(p.Options.AppDirectory),
		slogecho.New(p.Logger),
    )

	return Router{e}
}

func (r Router) MapHome(home echo.HandlerFunc, mw ...echo.MiddlewareFunc) {
    r.GET("/", home, mw...)
}

func (r Router) MapLinode(linode *controller.Linode, mw ...echo.MiddlewareFunc) {
	route := r.Group("/linode")
	route.GET("/sse", linode.SSE)
	route.POST("/boot", linode.Boot, mw...)
	route.POST("/reboot", linode.Reboot, mw...)
	route.POST("/shutdown", linode.Shutdown, mw...)
}

func (r Router) MapSystemd(systemd *controller.Systemd, mw ...echo.MiddlewareFunc) {
	route := r.Group("/systemd/:instance")
	route.GET("/sse", systemd.SSE)
	route.POST("/enable", systemd.Enable, mw...)
	route.POST("/disable", systemd.Disable, mw...)
	route.POST("/start", systemd.Start, mw...)
	route.POST("/stop", systemd.Stop, mw...)
	route.POST("/restart", systemd.Restart, mw...)
}

func (r Router) Start(addr string) error {
	go func() {
		_ = r.Echo.Start(addr)
	}()

	return nil
}
