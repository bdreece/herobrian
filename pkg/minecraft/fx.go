package minecraft

import (
	"github.com/bdreece/herobrian/pkg/minecraft/host"
	"github.com/bdreece/herobrian/pkg/minecraft/host/ec2"
	"github.com/bdreece/herobrian/pkg/minecraft/instance/ssh"
	"go.uber.org/fx"
)

var Module = fx.Module("minecraft",
	fx.Provide(
		fx.Annotate(
			ec2.NewProvider,
			fx.As(fx.Self()),
			fx.As(new(host.Provider)),
		),
		ssh.NewProviderFactory,
	),
)
