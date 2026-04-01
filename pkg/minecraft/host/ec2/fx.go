package ec2

import (
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/pkg/minecraft/host"
)

var Module = fx.Module("ec2",
	fx.Provide(
		fx.Annotate(
			NewProvider,
			fx.As(fx.Self()),
			fx.As(new(host.Provider)),
			fx.As(new(host.Describer)),
			fx.As(new(host.Starter)),
			fx.As(new(host.Stopper)),
			fx.As(new(host.Restarter)),
			fx.As(new(host.Checker)),
		),
	),
)
