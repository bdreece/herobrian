package herobrian

import (
	"context"
	"fmt"
	"maps"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	echorenderer "github.com/bdreece/echo-renderer"
	echovalidator "github.com/bdreece/echo-validator"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"go.uber.org/config"
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/internal/controller"
	"github.com/bdreece/herobrian/internal/logger"
	"github.com/bdreece/herobrian/internal/router"
	"github.com/bdreece/herobrian/pkg/auth"
	"github.com/bdreece/herobrian/pkg/identity"
	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/bdreece/herobrian/pkg/systemd"
	"github.com/bdreece/herobrian/web"
)

type Args struct {
	Port       int
	ConfigPath string
}

func (args Args) Addr() string {
	return fmt.Sprintf(":%d", args.Port)
}

func New(args Args) *fx.App {
	return fx.New(
		fx.Supply(args),
		fx.Provide(func(args Args) (config.Provider, error) {
			f, err := os.Open(args.ConfigPath)
			if err != nil {
				return nil, fmt.Errorf("failed to open config file %q: %v", args.ConfigPath, err)
			}

			return config.NewYAML(config.Source(f))
		}),
		fx.Provide(
			logger.Configure,
			logger.New,
		),
		fx.Provide(
			fx.Annotate(
				identity.NewSessionStore,
				fx.As(new(sessions.Store)),
			),
			fx.Annotate(
				identity.NewCookieAuthenticator,
				fx.As(new(identity.Authenticator)),
				fx.As(new(identity.SignInManager)),
			),
		),
		fx.Supply(&echorenderer.Options{
			FS:      web.Templates,
			Include: []string{"*.gotmpl"},
			Funcs: func(c echo.Context) template.FuncMap {
				funcs := make(template.FuncMap)

				maps.Copy(funcs, sprig.FuncMap())
				maps.Copy(funcs, template.FuncMap{
					"context": func() echo.Context { return c },
				})

				return funcs
			},
		}),
		fx.Supply(
			fx.Annotate(
				echovalidator.Default,
				fx.As(new(echo.Validator)),
			),
		),
		fx.Provide(
			fx.Annotate(
				echorenderer.New,
				fx.As(new(echo.Renderer)),
			),
		),
		fx.Provide(
			linode.Configure,
			linode.NewHTTP,
			linode.NewEmitter,
		),
		fx.Provide(
			systemd.ConfigureSSH,
			systemd.NewSSH,
			systemd.NewEmitter,
			systemd.NewServiceFactory,
		),
		fx.Provide(
			auth.Configure,
			auth.NewMiddleware,
		),
		fx.Provide(
			controller.Home,
			controller.NewLinode,
			controller.NewSystemd,
		),
		fx.Provide(
			router.Configure,
			router.New,
		),
		fx.Decorate(func(router router.Router, p struct {
			fx.In

			Home           echo.HandlerFunc
			Linode         *controller.Linode
			Systemd        *controller.Systemd
			AuthMiddleware echo.MiddlewareFunc

			Args      Args
			Lifecycle fx.Lifecycle
		}) router.Router {
			router.MapHome(p.Home, p.AuthMiddleware)
			router.MapLinode(p.Linode, p.AuthMiddleware)
			router.MapSystemd(p.Systemd, p.AuthMiddleware)

			p.Lifecycle.Append(fx.Hook{
				OnStart: func(context.Context) error {
					return router.Start(p.Args.Addr())
				},

				OnStop: func(ctx context.Context) error {
					return router.Shutdown(ctx)
				},
			})

			return router
		}),
		fx.Invoke(func(router.Router) {}),
	)
}
