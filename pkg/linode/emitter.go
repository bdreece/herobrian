//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Topic -trimprefix Topic
package linode

import (
	"context"
	"fmt"
	"time"

	"github.com/bdreece/herobrian/pkg/event"
	"github.com/bdreece/herobrian/pkg/worker"
)

type Emitter interface {
	event.Emitter[Topic, Status]
}

type emitter struct {
	event.Emitter[Topic, Status]

	worker worker.Service
}

func (e *emitter) Close() error {
	if err := e.Emitter.Close(); err != nil {
		return err
	}

	if err := e.worker.Stop(context.Background()); err != nil {
		return err
	}

	return nil
}

func NewEmitter(client Client) (Emitter, error) {
	e := event.NewEmitter[Topic, Status]()
	wrk := newWorkerService(workerParams{
		Client:   client,
		Emitter:  e,
		Interval: time.Second,
	})

	if err := wrk.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start linode worker: %w", err)
	}

	return &emitter{
		Emitter: e,
		worker:  wrk,
	}, nil
}
