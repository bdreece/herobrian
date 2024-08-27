package systemd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bdreece/herobrian/pkg/event"
	"github.com/bdreece/herobrian/pkg/worker"
)

type workerParams struct {
	Factory  *ServiceFactory
	Emitter  event.Emitter[string, Status]
	Instance string
	Interval time.Duration
	Logger   *slog.Logger
}

func newWorkerService(p workerParams) (worker.Service, error) {
	return worker.NewService(func(ctx context.Context) error {
		ticker := time.NewTicker(p.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				service, err := p.Factory.Create(p.Instance)
				if err != nil {
                    p.Logger.Error("something bad happened", slog.String("error", err.Error()))
					continue
				}

				status, err := service.Status(ctx)
				if err != nil {
                    p.Logger.Error("something bad happened", slog.String("error", err.Error()))
					return fmt.Errorf("failed to refresh service status: %w", err)
				}

				p.Emitter.Publish(p.Instance, *status)
			}
		}
	}), nil
}
