package systemd

import (
	"context"
	"errors"
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

func NewEmitter(services *ServiceFactory) (Emitter, error) {
	se := emitter{
		Emitter: event.NewEmitter[string, Status](),
	}

	errs := make([]error, 0)
	for _, unit := range services.Units() {
		wrk, err := newWorkerService(workerParams{
			Emitter:  &se,
            Factory: services,
            Instance: unit.Instance,
			Interval: time.Second,
		})

		if err != nil {
			errs = append(errs, err)
			continue
		}

		se.workers = append(se.workers, wrk)
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	if err := worker.Start(context.Background(), se.workers...); err != nil {
		return nil, err
	}

	return &se, nil
}
