package route

import (
	"net/url"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/spf13/viper"
)

func webServer() echo.MiddlewareFunc {
	env := viper.GetString("environment")
	skipper := func(c *echo.Context) bool {
		return strings.HasPrefix(c.Request().URL.Path, "/api")
	}

	if env == "production" {
		return middleware.StaticWithConfig(middleware.StaticConfig{
			Root:    viper.GetString("http:root_dir"),
			HTML5:   true,
			Skipper: skipper,
		})
	}

	url, _ := url.Parse(viper.GetString("http:proxy_url"))
	balancer := middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{URL: url},
	})

	return middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: balancer,
		Skipper:  skipper,
	})
}
