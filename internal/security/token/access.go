package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type AccessClaims struct {
	registeredClaims

	FirstName   string `json:"given_name"`
	LastName    string `json:"family_name"`
	DisplayName string `json:"preferred_username"`
	Picture     string `json:"picture,omitempty"`
	Role        string `json:"role,omitempty"`
}

type AccessHandler struct{ encoder[*AccessClaims] }

func NewAccessHandler() (*AccessHandler, error) {
	options := new(HandlerOptions)
	if err := viper.UnmarshalKey("jwt:access", options); err != nil {
		return nil, err
	}

	return &AccessHandler{encoder[*AccessClaims]{options}}, nil
}

func (h *AccessHandler) Decode(t string) (*AccessClaims, *jwt.Token, error) {
	token, err := jwt.ParseWithClaims(t, new(AccessClaims), h.keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return nil, nil, ErrInvalidClaims
	}

	return claims, token, nil
}

var _ Handler[*AccessClaims] = (*AccessHandler)(nil)
