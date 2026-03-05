package token

import (
	"go.uber.org/fx"
)

var Module = fx.Module("token",
	fx.Provide(
		asHandler(NewAccessHandler),
		asHandler(NewRefreshHandler),
		asHandler(NewInviteHandler),
	),
)

func asHandler[H Handler[C], C Claims](ctor func() (H, error)) any {
	return fx.Annotate(
		ctor,
		fx.As(fx.Self()),
		fx.As(new(Encoder[C])),
		fx.As(new(Decoder[C])),
		fx.As(new(Handler[C])),
	)
}
