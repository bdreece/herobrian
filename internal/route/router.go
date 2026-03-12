package route

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echovalidator "gopkg.in/bdreece/echo-validator.v1"
)

func NewRouter(controllers ...Controller) (*echo.Echo, error) {
	e := echo.New()
	e.HTTPErrorHandler = handleError
	e.Validator = echovalidator.Default

	e.Use(
		middleware.RequestLogger(),
		middleware.Recover(),
		middleware.ContextTimeout(5*time.Minute),
		middleware.BodyLimit(64*1024*1024),
		middleware.Secure(),
		webServer(),
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

	return e, nil
}
