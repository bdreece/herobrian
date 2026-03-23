package instance

import (
	"github.com/bdreece/herobrian/internal/route/host"
	"github.com/bdreece/herobrian/pkg/minecraft/instance"
	"github.com/labstack/echo/v5"
)

var contextKey string = "minecraft.host.instance"

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			host, _ := host.FromContext(c)
			id := c.Param("instance")
			provider, err := host.Instances(c.Request().Context())
			if err != nil {
				panic(err)
			}

			c.Set(contextKey, instance.NewClient(id, provider))

			return next(c)
		}
	}
}

func FromContext(c *echo.Context) (client *instance.Client, ok bool) {
	client, ok = c.Get(contextKey).(*instance.Client)

	return
}
