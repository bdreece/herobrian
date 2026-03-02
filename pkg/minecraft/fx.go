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
			NewHostHandler,
			fx.As(new(route.Router)),
			fx.ResultTags(`group:"routes"`),
		),
	),
)
