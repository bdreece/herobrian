package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/bdreece/herobrian/internal/middleware"
	"github.com/bdreece/herobrian/pkg/database"
	"github.com/bdreece/herobrian/pkg/identity"
	"github.com/bdreece/herobrian/pkg/token"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

type (
	Invite struct {
		db      database.Querier
		handler token.Handler[token.UserInviteClaims]
	}

	InviteParams struct {
		fx.In

		Querier database.Querier
		Handler token.Handler[token.UserInviteClaims]
	}
)

func (Invite) RenderSendInvite(c echo.Context) error {
	claims, ok := c.Get(middleware.ClaimsContextKey).(*identity.ClaimSet)
	if !ok || claims == nil {
		return fmt.Errorf("failed to get claims from request context")
	}

	roles := make(map[identity.Role]string)
	if claims.Role >= 0 {
		roles[identity.RoleUser] = "Users can access the system controls"
	}
	if claims.Role >= 1 {
		roles[identity.RoleModerator] = "Moderators can invite other users"
	}
	if claims.Role >= 2 {
		roles[identity.RoleAdmin] = "Admins can do other things, I haven't gotten that far yet"
	}
	if claims.Role >= 3 {
		roles[identity.RoleSuper] = "Fuck you, eat shit"
	}

	return c.Render(http.StatusOK, "send-invite.gotmpl", echo.Map{
		"Roles": roles,
	})
}

func (controller *Invite) SendInvite(c echo.Context) error {
	model := new(struct {
		Role int `form:"role" validate:"min=0,max=3"`
	})
	if err := c.Bind(model); err != nil {
		return err
	}
	if err := c.Validate(model); err != nil {
		return err
	}

	t, err := controller.handler.Sign(&token.UserInviteClaims{
		RoleID: model.Role,
	})
	if err != nil {
		return err
	}

	var proto string
	if c.IsTLS() {
		proto = "https"
	} else {
		proto = "http"
	}
	return c.HTML(http.StatusOK, fmt.Sprintf(`
        <div class="bg-neutral-200 rounded p-2 overflow-x-scroll col-span-2">
            <strong class="font-bold">Invite URL:</strong>
            <input
                class="input"
                type="text"
                value="%s://%s/invite/%s"
                is="invite-link"
            >
        </div>
    `, proto, c.Request().Host, t))
}

func (controller *Invite) RenderAcceptInvite(c echo.Context) error {
	model := new(struct {
		Token string `param:"token" validate:"required"`
	})
	if err := c.Bind(model); err != nil {
		return err
	}
	if err := c.Validate(model); err != nil {
		return err
	}

	_, err := controller.handler.Verify(model.Token)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "accept-invite.gotmpl", echo.Map{
		"Token": model.Token,
	})
}

func (controller *Invite) AcceptInvite(c echo.Context) error {
	model := new(struct {
		Token    string `param:"token" validate:"required"`
		Username string `form:"username" validate:"required,max=127"`
		Password string `form:"password" validate:"required,min=8,max=127"`
		_        string `form:"confirmPassword" validate:"required,min=8,max=127,eqfield=Password"`
	})
	if err := c.Bind(model); err != nil {
		return err
	}
	if err := c.Validate(model); err != nil {
		return err
	}

	claims, err := controller.handler.Verify(model.Token)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = controller.db.CreateUser(c.Request().Context(), database.CreateUserParams{
		Username:     model.Username,
		PasswordHash: base64.StdEncoding.EncodeToString(hash),
		RoleID:       int64(claims.RoleID),
	})
	if err != nil {
		return err
	}

	c.Response().Header().Add("HX-Location", "/login")
	return c.NoContent(http.StatusOK)
}

func NewInvite(p InviteParams) *Invite {
	return &Invite{p.Querier, p.Handler}
}
