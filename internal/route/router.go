package route

import (
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/spf13/viper"
	echovalidator "gopkg.in/bdreece/echo-validator.v1"
)

func NewRouter(controllers ...Controller) (*echo.Echo, error) {
	e := echo.New()
	e.Validator = echovalidator.Default

	e.Use(
		middleware.RequestLogger(),
		middleware.Recover(),
		middleware.ContextTimeout(5*time.Minute),
		middleware.BodyLimit(64*1024*1024),
		middleware.Secure(),
	)

	api := e.Group("/api")
	for _, controller := range controllers {
		for _, route := range controller.Routes() {
			info, err := api.AddRoute(route)
			if err != nil {
				return nil, err
			}

			slog.Debug("added route", "route", info)
		}
	}

	var appServer echo.MiddlewareFunc

	env := viper.GetString("environment")
	skipper := func(c *echo.Context) bool {
		return strings.HasPrefix(c.Request().URL.Path, "/api")
	}

	if env == "production" {
		appServer = middleware.StaticWithConfig(middleware.StaticConfig{
			Root:    viper.GetString("app:root_dir"),
			HTML5:   true,
			Skipper: skipper,
		})
	} else {
		url, _ := url.Parse(viper.GetString("app:proxy_url"))
		balancer := middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: url},
		})

		appServer = middleware.ProxyWithConfig(middleware.ProxyConfig{
			Balancer: balancer,
			Skipper:  skipper,
		})
	}

	e.Use(appServer)

	return e, nil
}
