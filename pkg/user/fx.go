package user

import (
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/internal/route"
)

var Module = fx.Module("user",
	fx.Provide(
		fx.Annotate(
			NewController,
			fx.As(new(route.Controller)),
			fx.ResultTags(`group:"controllers"`),
		),
	),
)
