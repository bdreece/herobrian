package instance

import (
	"github.com/bdreece/herobrian/internal/route"
	"go.uber.org/fx"
)

var Module = fx.Module("instance", route.ProvideController(NewController))
