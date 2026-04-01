package user

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/bdreece/herobrian/internal/database"
	"github.com/bdreece/herobrian/internal/security/token"
	"github.com/golang-jwt/jwt/v5"
)

func newAccessClaims(user *database.User) *token.AccessClaims {
	claims := token.AccessClaims{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DisplayName: user.DisplayName,
		Picture:     user.PictureURL,
	}

	now := time.Now()
	claims.Subject = fmt.Sprint(user.ID)
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)

	return &claims
}

func newRefreshClaims(user *database.User, accessToken string) *token.RefreshClaims {
	hash := md5.Sum([]byte(accessToken))
	claims := token.RefreshClaims{
		ATHash: base64.StdEncoding.EncodeToString(hash[:]),
	}

	now := time.Now()
	claims.Subject = fmt.Sprint(user.ID)
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)

	return &claims
}
