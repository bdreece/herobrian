package identity

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/go-crypt/crypt/algorithm"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewLoginHandler(db *database.Queries, signer *TokenSigner, decoder algorithm.Decoder) echo.HandlerFunc {
	return func(c echo.Context) error {
		var form struct {
			DisplayName string `form:"displayName"`
			Password    string `form:"password"`
			RememberMe  bool   `form:"rememberMe"`
		}

		if err := c.Bind(&form); err != nil {
			return echo.ErrBadRequest.WithInternal(err)
		}

		if err := c.Validate(&form); err != nil {
			return echo.ErrBadRequest.WithInternal(err)
		}

		ctx := c.Request().Context()

		user, err := db.FindUserByDisplayName(ctx, database.FindUserByDisplayNameParams{
			DisplayName: form.DisplayName,
		})
		if err != nil && err != sql.ErrNoRows {
			return echo.ErrBadGateway.WithInternal(err)
		} else if err == sql.ErrNoRows {
			return echo.ErrUnauthorized.WithInternal(err)
		}

		digest, _ := decoder.Decode(user.PasswordHash)
		if !digest.Match(form.Password) {
			return echo.ErrUnauthorized.WithInternal(err)
		}

		now := time.Now()
		claims := Claims{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			DisplayName: user.DisplayName,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   fmt.Sprint(user.ID),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
		}

		jwt, _, err := signer.Sign(&claims)
		if err != nil {
			return echo.ErrInternalServerError.WithInternal(err)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"accessToken": jwt,
		})
	}
}
