package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/bdreece/herobrian/internal/middleware"
	"github.com/bdreece/herobrian/pkg/database"
	"github.com/bdreece/herobrian/pkg/email"
	"github.com/bdreece/herobrian/pkg/identity"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

type (
	Auth struct {
		db      database.Querier
		client  email.Client
		mgr     identity.SignInManager
	}

	AuthParams struct {
		fx.In

		Querier       database.Querier
		EmailClient   email.Client
		SignInManager identity.SignInManager
	}

	authLoginModel struct {
		Username string `form:"username" validate:"required,max=127"`
		Password string `form:"password" validate:"required,min=8,max=127"`
	}
)

func (Auth) RenderLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.gotmpl", echo.Map{})
}

func (controller *Auth) Login(c echo.Context) error {
	if claims, ok := c.Get(middleware.ClaimsContextKey).(*identity.ClaimSet); ok && claims != nil {
		// user already logged in
		return c.Redirect(http.StatusFound, "/")
	}

	// parse form
	model := new(authLoginModel)
	if err := c.Bind(model); err != nil {
		return err
	}
	if err := c.Validate(model); err != nil {
		return err
	}

	// find user
	user, err := controller.db.FindUserByUsername(c.Request().Context(), model.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// compare hash
	hash, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to decode user password hash: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword(hash, []byte(model.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// create session
	err = controller.mgr.SignIn(c, &identity.ClaimSet{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.RoleID,
	})
	if err != nil {
		return fmt.Errorf("failed to sign in user: %w", err)
	}

	c.Response().Header().Add("HX-Location", "/")
	return c.NoContent(http.StatusOK)
}

func (controller *Auth) RenderLogout(c echo.Context) error {
	if err := controller.mgr.SignOut(c); err != nil {
		return err
	}

	c.Response().Header().Add("HX-Location", "/login")
	return c.NoContent(http.StatusOK)
}

func NewAuth(p AuthParams) *Auth {
	return &Auth{p.Querier, p.EmailClient, p.SignInManager}
}
