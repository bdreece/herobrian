package instance

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/bdreece/herobrian/internal/route"
)

type Controller struct{}

func NewController() *Controller {
	return new(Controller)
}

func (i *Controller) Routes() []echo.Route {
	return []echo.Route{
		{Method: http.MethodGet, Path: "/host/:hostname/instance", Handler: i.List},
	}
}

var _ route.Controller = (*Controller)(nil)
