package host

import (
	"github.com/bdreece/herobrian/internal/route"
	"go.uber.org/fx"
)

var Module = fx.Module("host",
	route.ProvideController(NewController),
)
