package user

import (
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/internal/route"
)

var Module = fx.Module("user",
	fx.Provide(
		fx.Annotate(
			NewHandler,
			fx.As(new(route.Router)),
			fx.ResultTags(`group:"routes"`),
		),
	),
)
