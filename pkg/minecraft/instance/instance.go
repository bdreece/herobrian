package instance

import (
	"context"
	"io"
	"net"
)

type Info struct {
	Name    string `json:"name"`
	Server  string `json:"server"`
	Flavor  string `json:"flavor"`
	Version string `json:"version"`
	IconURL string `json:"iconUrl"`
}

type InstanceProvider interface {
	Instances(ctx context.Context, ids ...string) ([]Info, error)
	StartInstance(ctx context.Context, id string) error
	StopInstance(ctx context.Context, id string) error
	RestartInstance(ctx context.Context, id string) error
	TraceInstance(ctx context.Context, id string) (*Tracer, error)
	DialInstance(ctx context.Context, id string) (InstanceConn, error)
}

type InstanceConn interface {
	io.Closer
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Execute(cmd string) (string, error)
}
