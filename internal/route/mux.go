package route

import (
	"log/slog"
	"net/url"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/spf13/viper"
	echovalidator "gopkg.in/bdreece/echo-validator.v1"
)

func NewMux(routers ...Router) (*echo.Echo, error) {
	e := echo.New()
	e.Validator = echovalidator.Default

	e.Use(
		middleware.RequestLogger(),
		middleware.Recover(),
	)

	api := e.Group("/api")
	for _, router := range routers {
		for _, r := range router.Routes() {
			info, err := api.AddRoute(r)
			if err != nil {
				return nil, err
			}

			slog.Debug("added route", "route", info)
		}
	}

	env := viper.GetString("environment")
	var appServer echo.MiddlewareFunc
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
