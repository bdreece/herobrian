package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type InviteClaims struct {
	registeredClaims

	DisplayName string `json:"preferred_username"`
	Role        string `json:"role"`
}

type InviteHandler struct{ encoder[*InviteClaims] }

func NewInviteHandler() (*InviteHandler, error) {
	options := new(HandlerOptions)
	if err := viper.UnmarshalKey("jwt:invite", options); err != nil {
		return nil, err
	}

	return &InviteHandler{encoder[*InviteClaims]{options}}, nil
}

func (h *InviteHandler) Decode(t string) (*InviteClaims, *jwt.Token, error) {
	token, err := jwt.ParseWithClaims(t, new(InviteClaims), h.keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*InviteClaims)
	if !ok {
		return nil, nil, ErrInvalidClaims
	}

	return claims, token, nil
}

var _ Handler[*InviteClaims] = (*InviteHandler)(nil)
