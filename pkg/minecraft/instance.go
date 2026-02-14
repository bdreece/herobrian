package minecraft

import (
	"context"
	"io"
	"net"
)

type InstanceInfo struct {
}

type InstanceTrace struct {
}

type InstanceTracer struct {
	traces chan InstanceTrace
	cancel context.CancelFunc
}

// Close implements [io.Closer].
func (tracer *InstanceTracer) Close() error {
	tracer.cancel()
	return nil
}

func (tracer *InstanceTracer) Traces() <-chan InstanceTrace {
	return tracer.traces
}

type InstanceProvider interface {
	Instances(ctx context.Context, ids ...string) ([]InstanceInfo, error)
	StartInstance(ctx context.Context, id string) error
	StopInstance(ctx context.Context, id string) error
	RestartInstance(ctx context.Context, id string) error
	TraceInstance(ctx context.Context, id string) (*InstanceTracer, error)
	DialInstance(ctx context.Context, id string) (InstanceConn, error)
}

type InstanceConn interface {
	io.Closer
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Execute(cmd string) (string, error)
}
