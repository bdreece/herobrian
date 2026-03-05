package user

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/security/token"
)

type loginForm struct {
	DisplayName string `form:"displayName" validate:"required"`
	Password    string `form:"password"    validate:"required"`
	RememberMe  *bool  `form:"rememberMe"`
}

type loginResult struct {
	AccessToken string `json:"accessToken"`
}

func (self *Controller) Login(c *echo.Context) error {
	var form loginForm

	if err := c.Bind(&form); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	if err := c.Validate(&form); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	self.logger.Debug("parsed form", "form", form)

	ctx := c.Request().Context()
	user, err := self.querier.FindUserByDisplayName(ctx, database.FindUserByDisplayNameParams{
		DisplayName: form.DisplayName,
	})

	if err != nil && err != sql.ErrNoRows {
		return echo.ErrBadGateway.Wrap(err)
	} else if err == sql.ErrNoRows {
		return echo.ErrUnauthorized.Wrap(err)
	}

	self.logger.Debug("queried user", "user", user.ID)
	digest, _ := self.passwordHasher.Decode(user.PasswordHash)
	if !digest.Match(form.Password) {
		return echo.ErrUnauthorized.Wrap(err)
	}

	self.logger.Debug("authenticated user")

	now := time.Now()
	accessClaims := token.AccessClaims{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DisplayName: user.DisplayName,
	}

	accessClaims.SetSubject(fmt.Sprint(user.ID))
	accessClaims.SetNotBefore(jwt.NewNumericDate(now))
	accessClaims.SetIssuedAt(jwt.NewNumericDate(now))

	accessToken, _, err := self.accessTokenEncoder.Encode(&accessClaims)
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	self.logger.Debug("signed user access token")

	if form.RememberMe != nil && *form.RememberMe {
		atHash := md5.Sum([]byte(accessToken))
		refreshClaims := token.RefreshClaims{
			ATHash: base64.StdEncoding.EncodeToString(atHash[:]),
		}

		refreshClaims.SetSubject(fmt.Sprint(user.ID))
		refreshClaims.SetNotBefore(jwt.NewNumericDate(now))
		refreshClaims.SetIssuedAt(jwt.NewNumericDate(now))

		refreshToken, _, err := self.refreshTokenEncoder.Encode(&refreshClaims)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		var cookie http.Cookie
		if err := viper.UnmarshalKey("http:cookie:refresh", &cookie); err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		cookie.Value = refreshToken

		c.SetCookie(&cookie)
	}

	return c.JSON(http.StatusOK, loginResult{
		AccessToken: accessToken,
	})
}
