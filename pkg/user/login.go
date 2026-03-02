package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/identity"
)

type loginForm struct {
	DisplayName string `form:"displayName"`
	Password    string `form:"password"`
	RememberMe  bool   `form:"rememberMe"`
}

type loginResult struct {
	AccessToken string `json:"accessToken"`
}

func (h *Handler) Login(c *echo.Context) error {
	var form loginForm

	if err := c.Bind(&form); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	if err := c.Validate(&form); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	ctx := c.Request().Context()
	user, err := h.queries.FindUserByDisplayName(ctx, database.FindUserByDisplayNameParams{
		DisplayName: form.DisplayName,
	})

	if err != nil && err != sql.ErrNoRows {
		return echo.ErrBadGateway.Wrap(err)
	} else if err == sql.ErrNoRows {
		return echo.ErrUnauthorized.Wrap(err)
	}

	digest, _ := h.hasher.Decode(user.PasswordHash)
	if !digest.Match(form.Password) {
		return echo.ErrUnauthorized.Wrap(err)
	}

	now := time.Now()
	claims := identity.Claims{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DisplayName: user.DisplayName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(user.ID),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	jwt, _, err := h.signer.Sign(&claims)
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	return c.JSON(http.StatusOK, loginResult{
		AccessToken: jwt,
	})
}
