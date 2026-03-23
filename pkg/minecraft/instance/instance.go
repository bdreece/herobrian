package instance

import (
	"context"
	"io"
	"net"
	"time"
)

type Info struct {
	Name    string `json:"name"`
	Server  string `json:"server"`
	Flavor  string `json:"flavor"`
	Version string `json:"version"`
	IconURL string `json:"iconUrl"`
}

type Trace struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type Conn interface {
	io.Closer
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Execute(cmd string) (string, error)
}

type Dialer interface {
	DialInstance(ctx context.Context, id string) (Conn, error)
}

type Describer interface {
	DescribeInstance(ctx context.Context, id string) (map[string]string, error)
}

type Starter interface {
	StartInstance(ctx context.Context, id string) error
}

type Stopper interface {
	StopInstance(ctx context.Context, id string) error
}

type Restarter interface {
	RestartInstance(ctx context.Context, id string) error
}

type Tracer interface {
	TraceInstance(ctx context.Context, id string) (<-chan Trace, error)
}

type Provider interface {
	Dialer
	Describer
	Starter
	Stopper
	Restarter
	Tracer
}

type ProviderFactory interface {
	Instances(ctx context.Context, hostname string) (Provider, error)
}

type ProviderFactoryFunc func(ctx context.Context, hostname string) (Provider, error)

var _ ProviderFactory = ProviderFactoryFunc(nil)

func (fn ProviderFactoryFunc) Instances(ctx context.Context, hostname string) (Provider, error) {
	return fn(ctx, hostname)
}
