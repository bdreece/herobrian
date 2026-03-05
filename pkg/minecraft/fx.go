package minecraft

import (
	"github.com/bdreece/herobrian/internal/route"
	"go.uber.org/fx"
)

var Module = fx.Module("minecraft",
	fx.Provide(
		fx.Annotate(
			NewEC2Provider,
			fx.As(new(HostProvider)),
		),
		fx.Annotate(
			NewHostController,
			fx.As(new(route.Controller)),
			fx.ResultTags(`group:"controllers"`),
		),
	),
)
