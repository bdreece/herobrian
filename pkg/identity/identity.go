package identity

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	ClaimSet struct {
		ID        int64  `mapstructure:"id"`
		Username  string `mapstructure:"username"`
		FirstName string `mapstructure:"first_name"`
		LastName  string `mapstructure:"last_name"`
		Role      int64  `mapstructure:"role"`
	}

	Authenticator interface {
		Authenticate(echo.Context) (*ClaimSet, error)
	}

	AuthorizerFunc func(*ClaimSet) error

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

var (
	ErrUnauthenticated = errors.New("user is unauthenticated")
	ErrUnauthorized    = errors.New("user is unauthorized")
)

func (fn AuthorizerFunc) Authorize(claims *ClaimSet) error {
	return fn(claims)
}

func AuthorizeFunc(fn AuthorizerFunc) Authorizer { return fn }

var DefaultAuthorizer = AuthorizeFunc(func(cs *ClaimSet) error {
	if cs == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "no claims found")
	}

	return nil
})
