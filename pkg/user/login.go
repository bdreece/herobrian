package user

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"

	"github.com/bdreece/herobrian/internal/database"
)

type loginForm struct {
	DisplayName string `form:"displayName" validate:"required"`
	Password    string `form:"password"    validate:"required"`
	RememberMe  *bool  `form:"rememberMe"`
}

type loginResult struct {
	AccessToken string `json:"accessToken"`
}

func (u *Controller) Login(c *echo.Context) error {
	var form loginForm

	if err := c.Bind(&form); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	if err := c.Validate(&form); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	u.logger.Debug("parsed form", "form", form)

	ctx := c.Request().Context()
	user, err := u.querier.FindUserByDisplayName(ctx, database.FindUserByDisplayNameParams{
		DisplayName: form.DisplayName,
	})

	if err != nil && err != sql.ErrNoRows {
		return echo.ErrBadGateway.Wrap(err)
	} else if err == sql.ErrNoRows {
		return echo.ErrUnauthorized.Wrap(err)
	}

	u.logger.Debug("queried user", "user", user.ID)
	digest, _ := u.passwordHasher.Decode(user.PasswordHash)
	if !digest.Match(form.Password) {
		return echo.ErrUnauthorized.Wrap(err)
	}

	u.logger.Debug("authenticated user")

	accessToken, _, err := u.accessTokenEncoder.Encode(newAccessClaims(user))
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	u.logger.Debug("signed user access token")

	if form.RememberMe != nil && *form.RememberMe {
		refreshToken, _, err := u.refreshTokenHandler.Encode(newRefreshClaims(user, accessToken))
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		cookie := http.Cookie{
			Name:     viper.GetString("http:cookie:refresh"),
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   604800,
		}

		c.SetCookie(&cookie)
	}

	return c.JSON(http.StatusOK, loginResult{
		AccessToken: accessToken,
	})
}
