package identity

import (
	"context"

	"github.com/labstack/echo/v4"
)

type (
	ClaimSet struct {
		ID        int32  `mapstructure:"id"`
		UserName  string `mapstructure:"username"`
		FirstName string `mapstructure:"first_name"`
		LastName  string `mapstructure:"last_name"`
	}

	Authenticator interface {
		Authenticate(echo.Context) (*ClaimSet, error)
	}

	Authorizer interface {
		Authorize(*ClaimSet) error
	}

	SignInManager interface {
		SignIn(echo.Context, *ClaimSet) error
		SignOut(echo.Context) error
	}

	PasswordManager interface {
		SetPassword(ctx context.Context, claims *ClaimSet, oldPassword, newPassword string) error
		SendPasswordReset(ctx context.Context, email string) error
		ConfirmPasswordReset(ctx context.Context, email, password, confirmation string) error
	}
)
