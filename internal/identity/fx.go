package identity

import "go.uber.org/fx"

var Module = fx.Module("identity",
	fx.Provide(
		fx.Annotate(
			newAccessTokenSigner,
			fx.ResultTags(`name:"access"`),
		),
		fx.Annotate(
			newInviteTokenSigner,
			fx.ResultTags(`name:"invite"`),
		),
		fx.Annotate(
			NewArgon2,
			fx.As(new(PasswordHasher)),
		),
	),
)

func newAccessTokenSigner() (*TokenSigner, error) {
	return NewTokenSigner(AccessToken)
}

func newInviteTokenSigner() (*TokenSigner, error) {
	return NewTokenSigner(InviteToken)
}
