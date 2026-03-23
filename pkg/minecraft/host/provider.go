package host

import (
	"context"

	"github.com/bdreece/herobrian/pkg/minecraft/instance"
)

type Checker interface {
	CheckHosts(ctx context.Context, ids ...string) (map[string]Status, error)
}

type Describer interface {
	DescribeHosts(ctx context.Context, ids ...string) (map[string]*Info, error)
}

type Starter interface {
	StartHosts(ctx context.Context, ids ...string) error
}

type Stopper interface {
	StopHosts(ctx context.Context, ids ...string) error
}

type Restarter interface {
	RestartHosts(ctx context.Context, ids ...string) error
}

type Provider interface {
	Checker
	Describer
	Starter
	Stopper
	Restarter
	instance.ProviderFactory
}
