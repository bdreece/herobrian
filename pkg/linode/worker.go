package linode

import (
	"context"
	"time"

	"github.com/bdreece/herobrian/pkg/event"
	"github.com/bdreece/herobrian/pkg/worker"
)

type WorkerOptions struct {
	Interval time.Duration `yaml:"interval"`
}

type workerParams struct {
	Client   Client
	Emitter  event.Emitter[Topic, Status]
	Interval time.Duration
}

func newWorkerService(p workerParams) worker.Service {
	return worker.NewService(func(ctx context.Context) error {
		ticker := time.NewTicker(p.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				if !p.Emitter.Ready() {
					continue
				}

				status, err := p.Client.InstanceStatus(ctx)
				if err != nil {
					return err
				}

				p.Emitter.Publish(TopicStatus, *status)
			}
		}
	})
}
