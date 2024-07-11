package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/bdreece/herobrian/pkg/linode"
	"github.com/labstack/echo/v4"
)

const runningEvent string = "event: running\n" +
    "data: <button class=\"rounded-full bg-primary\" hx-post=\"/shutdown\" hx-swap=\"outerHTML\">Shutdown</button>\n\n"
const offlineEvent string = "event: offline\n" +
    "data: <button class=\"rounded-full bg-primary\" hx-post=\"/boot\" hx-swap=\"outerHTML\">Boot</button>\n\n"

func SSE(client *linode.Client) echo.HandlerFunc {
    return func(c echo.Context) error {
        w := c.Response()
        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")

        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-c.Request().Context().Done():
                return nil
            case <-ticker.C:
                var buf bytes.Buffer
                status, err := client.InstanceStatus(c.Request().Context())
                if err != nil {
                    return echo.NewHTTPError(http.StatusFailedDependency, err.Error())
                }

                _, _ = buf.WriteString(fmt.Sprintf("event: status\ndata: %s\n\n", status))
                if status == "running" {
                    _, _ = buf.WriteString(runningEvent)
                } else if status == "offline" {
                    _, _ = buf.WriteString(offlineEvent)
                }

                w.Write(buf.Bytes())
                w.Flush()
            }
        }
    }
}
