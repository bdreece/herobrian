package router

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/internal/controller"
	mw "github.com/bdreece/herobrian/internal/middleware"
	"github.com/bdreece/herobrian/pkg/identity"
)

type Params struct {
	fx.In

	Renderer      echo.Renderer
	Validator     echo.Validator
	Logger        *slog.Logger
	SessionStore  sessions.Store
	Authenticator identity.Authenticator
	Options       *Options
}

type Router struct {
	*echo.Echo

	authenticate   echo.MiddlewareFunc
	authorize      echo.MiddlewareFunc
	allowModerator echo.MiddlewareFunc
}

func (r Router) MapHome(home *controller.Home) {
	r.GET("/", home.RenderIndex, r.authenticate, r.authorize)
}

func (r Router) MapAuth(auth *controller.Auth) {
	r.GET("/login", auth.RenderLogin)
	r.POST("/login", auth.Login)
	r.GET("/logout", auth.RenderLogout)
}

func (r Router) MapInvite(invite *controller.Invite) {
	g := r.Group("/invite")
	g.GET("", invite.RenderSendInvite, r.authenticate, r.authorize, r.allowModerator)
	g.POST("", invite.SendInvite, r.authenticate, r.authorize, r.allowModerator)
	g.GET("/:token", invite.RenderAcceptInvite)
	g.POST("/:token", invite.AcceptInvite)
}

func (r Router) MapLinode(linode *controller.Linode) {
	route := r.Group("/linode", r.authenticate, r.authorize)
	route.GET("/sse", linode.SSE)
	route.POST("/boot", linode.Boot)
	route.POST("/reboot", linode.Reboot)
	route.POST("/shutdown", linode.Shutdown)
}

func (r Router) MapSystemd(systemd *controller.Systemd) {
	route := r.Group("/systemd/:instance", r.authenticate, r.authorize)
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

func New(p Params) Router {
	e := echo.New()
	e.Renderer = p.Renderer
	e.Validator = p.Validator
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		var location string
		switch code {
		case http.StatusUnauthorized:
			location = "/login"
		case http.StatusForbidden:
			location = ".."
		default:
			_ = c.NoContent(code)
			return
		}

		var model struct {
			HXRequest bool `header:"HX-Request"`
		}

		_ = c.Bind(&model)

		if model.HXRequest {
			c.Response().Header().Add("HX-Location", location)
			_ = c.NoContent(code)
		} else {
			_ = c.Redirect(http.StatusFound, location)
		}

	}

	e.Use(
		middleware.BodyLimit("4M"),
		middleware.Decompress(),
		middleware.Gzip(),
		middleware.CSRF(),
		middleware.Secure(),
		middleware.Static(p.Options.StaticDirectory),
		middleware.Static(p.Options.AppDirectory),
		session.Middleware(p.SessionStore),
		slogecho.New(p.Logger),
	)

	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "not-found.gotmpl", echo.Map{})
	})

	return Router{
		Echo:         e,
		authenticate: mw.Authenticate(p.Authenticator),
		authorize:    mw.Authorize(identity.DefaultAuthorizer),
		allowModerator: mw.Authorize(
			identity.DefaultAuthorizer,
			identity.RoleModerator,
		),
	}
}
