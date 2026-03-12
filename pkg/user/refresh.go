package user

import (
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
)

type refreshResult struct {
	AccessToken string `json:"accessToken"`
}

func (self *Controller) Refresh(c *echo.Context) error {
	cookie, err := c.Cookie(viper.GetString("http:cookie:refresh"))
	if err != nil {
		return echo.ErrUnauthorized.Wrap(err)
	}

	claims, _, err := self.refreshTokenHandler.Decode(cookie.Value)
	if err != nil {
		return echo.ErrUnauthorized.Wrap(err)
	}

	id, _ := strconv.ParseInt(claims.Subject, 10, 64)

	header := c.Request().Header.Get("Authorization")
	jwt := strings.TrimPrefix(header, "Bearer ")
	hash := md5.Sum([]byte(jwt))

	if claims.ATHash != base64.StdEncoding.EncodeToString(hash[:]) {
		return echo.ErrUnauthorized.Wrap(err)
	}

	params := database.FindUserByIDParams{
		ID: id,
	}

	user, err := self.querier.FindUserByID(c.Request().Context(), params)
	if err != nil {
		return echo.ErrBadGateway.Wrap(err)
	}

	accessToken, _, err := self.accessTokenEncoder.Encode(newAccessClaims(user))
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	refreshToken, _, err := self.refreshTokenHandler.Encode(newRefreshClaims(user, accessToken))
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	c.SetCookie(&http.Cookie{
		Name:     viper.GetString("http:cookie:refresh"),
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800,
	})

	return c.JSON(http.StatusOK, refreshResult{
		AccessToken: accessToken,
	})
}
