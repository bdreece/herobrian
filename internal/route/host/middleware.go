package host

import (
	"github.com/bdreece/herobrian/pkg/minecraft/host"
	"github.com/labstack/echo/v5"
)

var contextKey string = "minecraft.host"

func Middleware(provider host.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			hostname := c.Param("host")
			client, err := host.NewClient(hostname, provider)
			if err != nil {
				panic(err)
			}

			c.Set(contextKey, client)

			return next(c)
		}
	}
}

func FromContext(c *echo.Context) (client *host.Client, ok bool) {
	client, ok = c.Get(contextKey).(*host.Client)
	return
}
