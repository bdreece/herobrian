package identity

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/config"
)

type (
	CookieOptions struct {
		Name     string        `yaml:"name"`
		Path     string        `yaml:"path"`
		Domain   string        `yaml:"domain"`
		MaxAge   int           `yaml:"max_age"`
		Secure   bool          `yaml:"secure"`
		HttpOnly bool          `yaml:"http_only"`
		SameSite http.SameSite `yaml:"same_site"`
	}

	CookieAuthenticator struct {
		opts *CookieOptions
	}
)

func (opts CookieOptions) SessionOptions() *sessions.Options {
	return &sessions.Options{
		Path:     opts.Path,
		Domain:   opts.Domain,
		MaxAge:   opts.MaxAge,
		Secure:   opts.Secure,
		HttpOnly: opts.HttpOnly,
		SameSite: opts.SameSite,
	}
}

// Authenticate implements Authenticator.
func (ca *CookieAuthenticator) Authenticate(c echo.Context) (*ClaimSet, error) {
	sess, err := session.Get(ca.opts.Name, c)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to get session: %w", err),
			ErrUnauthenticated,
		)
	}

	if sess.IsNew {
		return nil, errors.Join(
			fmt.Errorf("session is new"),
			ErrUnauthenticated,
		)
	}

	claims := new(ClaimSet)
	if err = mapstructure.Decode(sess.Values, &claims); err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to decode claims: %w", err),
			ErrUnauthenticated,
		)
	}

	return claims, nil
}

// SignIn implements SignInManager.
func (ca *CookieAuthenticator) SignIn(c echo.Context, claims *ClaimSet) error {
	sess, err := session.Get(ca.opts.Name, c)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if sess.IsNew {
		sess.Options = ca.opts.SessionOptions()
	}

	if err = mapstructure.Decode(claims, &sess.Values); err != nil {
		return fmt.Errorf("failed to encode claims: %w", err)
	}

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

// SignOut implements SignInManager.
func (ca *CookieAuthenticator) SignOut(c echo.Context) error {
	sess, err := session.Get(ca.opts.Name, c)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if sess.IsNew {
		return nil
	}

	sess.Options.MaxAge = -1
	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func ConfigureCookie(provider config.Provider) (*CookieOptions, error) {
	opts := new(CookieOptions)
	if err := provider.Get("session.cookie").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure cookie options: %w", err)
	}

	return opts, nil
}

func NewCookieAuthenticator(opts *CookieOptions) *CookieAuthenticator {
	return &CookieAuthenticator{opts}
}

var (
	_ Authenticator = (*CookieAuthenticator)(nil)
	_ SignInManager = (*CookieAuthenticator)(nil)
)
