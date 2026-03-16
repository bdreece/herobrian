package minecraft

import (
	"github.com/bdreece/herobrian/internal/route"
	"github.com/bdreece/herobrian/pkg/minecraft/host"
	"go.uber.org/fx"
)

var Module = fx.Module("minecraft",
	fx.Provide(
		fx.Annotate(
			host.NewEC2Provider,
			fx.As(fx.Self()),
			fx.As(new(host.Provider)),
		),
		fx.Annotate(
			host.NewController,
			fx.As(new(route.Controller)),
			fx.ResultTags(`group:"controllers"`),
		),
	),
)
