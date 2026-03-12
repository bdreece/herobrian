package route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

func handleError(c *echo.Context, err error) {
	slog.Warn("error handling request", "error", err)

	if res, resErr := echo.UnwrapResponse(c.Response()); resErr == nil {
		if res.Committed {
			return
		}
	}

	var sc echo.HTTPStatusCoder

	code := http.StatusInternalServerError
	if errors.As(err, &sc) {
		if newCode := sc.StatusCode(); newCode != 0 {
			code = newCode
		}
	}

	if clientErr := c.NoContent(code); clientErr != nil {
		slog.Error("error writing response", "error", clientErr)
	}
}
