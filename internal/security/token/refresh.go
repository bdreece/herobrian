package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type RefreshClaims struct {
	registeredClaims

	ATHash string `json:"at_hash"`
}

type RefreshHandler struct{ encoder[*RefreshClaims] }

func NewRefreshHandler() (*RefreshHandler, error) {
	options := new(HandlerOptions)
	if err := viper.UnmarshalKey("jwt:invite", options); err != nil {
		return nil, err
	}

	return &RefreshHandler{encoder[*RefreshClaims]{options}}, nil
}

func (h *RefreshHandler) Decode(t string) (*RefreshClaims, *jwt.Token, error) {
	token, err := jwt.ParseWithClaims(t, new(RefreshClaims), h.keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return nil, nil, ErrInvalidClaims
	}

	return claims, token, nil
}

var _ Handler[*RefreshClaims] = (*RefreshHandler)(nil)
