package instance

import "context"

type Trace struct {
}

type Tracer struct {
	traces chan Trace
	cancel context.CancelFunc
}

// Close implements [io.Closer].
func (tracer *Tracer) Close() error {
	tracer.cancel()
	return nil
}

func (tracer *Tracer) Traces() <-chan Trace {
	return tracer.traces
}
