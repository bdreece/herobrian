package systemd

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/bdreece/herobrian/pkg/event"
	"github.com/bdreece/herobrian/pkg/worker"
)

type Emitter interface {
	event.Emitter[string, Status]
}

type emitter struct {
	event.Emitter[string, Status]

	workers []worker.Service
}

func (se *emitter) Close() error {
	if err := worker.Stop(context.Background(), se.workers...); err != nil {
		return err
	}

	if err := se.Emitter.Close(); err != nil {
		return err
	}

	return nil
}

func NewEmitter(services *ServiceFactory, logger *slog.Logger) (Emitter, error) {
	e := event.NewEmitter[string, Status]()

	errs := make([]error, 0)
	wrks := make([]worker.Service, 0, len(services.Units()))
	for _, unit := range services.Units() {
		wrk, err := newWorkerService(workerParams{
			Emitter:  e,
			Factory:  services,
			Instance: unit.Instance,
			Interval: 15 * time.Second,
			Logger:   logger,
		})

		if err != nil {
			errs = append(errs, err)
			continue
		}

		wrks = append(wrks, wrk)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	if err := worker.Start(context.Background(), wrks...); err != nil {
		return nil, err
	}

	return &emitter{
		Emitter: e,
		workers: wrks,
	}, nil
}
