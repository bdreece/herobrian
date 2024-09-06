package identity

import (
	"context"
	"github.com/bdreece/herobrian/pkg/email"
	"github.com/bdreece/herobrian/pkg/token"
	"go.uber.org/fx"
)

type (
	passwordManager struct {
		client  email.Client
		handler token.Handler[token.PasswordResetClaims]
	}

	PasswordManagerParams struct {
		fx.In

		EmailClient  email.Client
		TokenHandler token.Handler[token.PasswordResetClaims]
	}
)

// ConfirmPasswordReset implements PasswordManager.
func (p *passwordManager) ConfirmPasswordReset(ctx context.Context, email string, password string, confirmation string) error {
	panic("unimplemented")
}

// SendPasswordReset implements PasswordManager.
func (p *passwordManager) SendPasswordReset(ctx context.Context, email string) error {
	panic("unimplemented")
}

// SetPassword implements PasswordManager.
func (p *passwordManager) SetPassword(ctx context.Context, claims *ClaimSet, oldPassword string, newPassword string) error {
	panic("unimplemented")
}

func NewPasswordResetManager(p PasswordManagerParams) PasswordManager {
	return &passwordManager{p.EmailClient, p.TokenHandler}
}
