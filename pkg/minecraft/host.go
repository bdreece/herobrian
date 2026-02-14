package minecraft

import (
	"context"
	"errors"
)

var ErrHostNotFound = errors.New("minecraft: host not found")

type HostInfo struct {
	ID           string
	Architecture string
	DNSName      string
	Image        *ImageInfo
	IPAddress    string
	Platform     string
	Processor    ProcessorInfo
}

type ImageInfo struct {
	Type string
	// Memory size in MiB.
	Memory *int64
	// Storage space in GiB.
	Storage *int64
	Network string
}

type ProcessorInfo struct {
	CoreCount      *int32
	ThreadsPerCore *int32
}

type HostProvider interface {
	Hosts(ctx context.Context, ids ...string) ([]*HostInfo, error)
	StartHosts(ctx context.Context, ids ...string) error
	StopHosts(ctx context.Context, ids ...string) error
	RestartHosts(ctx context.Context, ids ...string) error
}
