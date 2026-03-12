package user

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/security"
	"github.com/bdreece/herobrian/internal/security/token"
)

type Controller struct {
	querier             database.Querier
	passwordHasher      security.PasswordHasher
	accessTokenEncoder  token.Encoder[*token.AccessClaims]
	inviteTokenEncoder  token.Encoder[*token.InviteClaims]
	refreshTokenHandler token.Handler[*token.RefreshClaims]
	logger              *slog.Logger
}

type ControllerParams struct {
	fx.In

	Querier             database.Querier
	PasswordHasher      security.PasswordHasher
	AccessTokenEncoder  token.Encoder[*token.AccessClaims]
	InviteTokenEncoder  token.Encoder[*token.InviteClaims]
	RefreshTokenHandler token.Handler[*token.RefreshClaims]
	Logger              *slog.Logger
}

func NewController(p ControllerParams) *Controller {
	return &Controller{
		p.Querier,
		p.PasswordHasher,
		p.AccessTokenEncoder,
		p.InviteTokenEncoder,
		p.RefreshTokenHandler,
		p.Logger.With("scope", "user.Controller"),
	}
}

func (self *Controller) Routes() []echo.Route {
	return []echo.Route{
		{Method: http.MethodPost, Path: "/user/activate", Handler: self.Activate},
		{Method: http.MethodPost, Path: "/user/invite", Handler: self.Invite},
		{Method: http.MethodPost, Path: "/user/login", Handler: self.Login},
		{Method: http.MethodPost, Path: "/user/logout", Handler: self.Logout},
		{Method: http.MethodPost, Path: "/user/refresh", Handler: self.Refresh},
	}
}
