package worker

import (
	"context"
	"errors"
	"time"
)

var (
	ErrWorkerStopped = errors.New("worker stopped")
	ErrWorkerTimeout = errors.New("worker failed to shutdown before timeout")
)

type RunFunc func(context.Context) error

type Service interface {
	// Start launches the routine's Run method on a new goroutine.
	Start(context.Context) error
	// Stop cancels the previously started goroutine.
	Stop(context.Context) error
}

type service struct {
	run    RunFunc
	cancel context.CancelCauseFunc
	err    error
}

func (ws *service) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	ctx, cancel := context.WithCancelCause(ctx)
	ws.cancel = cancel

	go func() {
		ws.err = ws.run(ctx)
	}()

	return nil
}

func (ws *service) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeoutCause(ctx, 5*time.Second, ErrWorkerTimeout)
	defer cancel()

	ch := make(chan struct{}, 1)
	go func() {
		ws.cancel(ErrWorkerStopped)
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return ws.err
	}
}

func NewService(run RunFunc) Service {
	return &service{run: run}
}

func Start(ctx context.Context, services ...Service) error {
	errch := make(chan error, len(services))
	for _, service := range services {
		go func(service Service) {
			errch <- service.Start(ctx)
		}(service)
	}

	errs := make([]error, 0)
	for range services {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errch:
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func Stop(ctx context.Context, services ...Service) error {
	errch := make(chan error, len(services))
	for _, service := range services {
		go func(service Service) {
			errch <- service.Stop(ctx)
		}(service)
	}

	errs := make([]error, 0)
	for range services {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errch:
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
