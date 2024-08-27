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
	route := r.Group("/linode", mw...)
	route.GET("/sse", linode.SSE)
	route.POST("/boot", linode.Boot)
	route.POST("/reboot", linode.Reboot)
	route.POST("/shutdown", linode.Shutdown)
}

func (r Router) MapSystemd(systemd *controller.Systemd, mw ...echo.MiddlewareFunc) {
	route := r.Group("/systemd/:instance", mw...)
	route.GET("/sse", systemd.SSE)
	route.POST("/enable", systemd.Enable)
	route.POST("/disable", systemd.Disable)
	route.POST("/start", systemd.Start)
	route.POST("/stop", systemd.Stop)
	route.POST("/restart", systemd.Restart)
}

func (r Router) Start(addr string) error {
	go func() {
		_ = r.Echo.Start(addr)
	}()

	return nil
}
