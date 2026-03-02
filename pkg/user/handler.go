package user

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/fx"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/identity"
)

type Handler struct {
	queries *database.Queries
	hasher  identity.PasswordHasher
	signer  *identity.TokenSigner
}

type HandlerParams struct {
	fx.In

	Queries *database.Queries
	Hasher  identity.PasswordHasher
	Signer  *identity.TokenSigner `name:"access"`
}

func NewHandler(p HandlerParams) *Handler {
	return &Handler{p.Queries, p.Hasher, p.Signer}
}

func (h *Handler) Routes() []echo.Route {
	return []echo.Route{
		{Method: http.MethodPost, Path: "/user/activate", Handler: h.Activate},
		{Method: http.MethodPost, Path: "/user/invite", Handler: h.Invite},
		{Method: http.MethodPost, Path: "/user/login", Handler: h.Login},
		{Method: http.MethodGet, Path: "/user/logout", Handler: h.Logout},
	}
}
