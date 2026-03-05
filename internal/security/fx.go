package security

import (
	"github.com/bdreece/herobrian/internal/security/token"
	"go.uber.org/fx"
)

var Module = fx.Module("security",
	token.Module,
	fx.Provide(
		fx.Annotate(
			NewArgon2,
			fx.As(new(PasswordHasher)),
		),
	),
)
