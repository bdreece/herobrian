package identity

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	jwt.RegisteredClaims
	FirstName   string `mapstructure:"given_name"`
	LastName    string `mapstructure:"family_name"`
	DisplayName string `mapstructure:"preferred_username"`
	Picture     string `mapstructure:"picture"`
	Role        string `mapstructure:"role"`
}

type TokenOptions struct {
	Audience string        `mapstructure:"aud"`
	Issuer   string        `mapstructure:"iss"`
	Secret   string        `mapstructure:"secret"`
	Lifetime time.Duration `mapstructure:"lifetime"`
}

type TokenSigner struct {
	options *TokenOptions
}

type TokenKind string

const (
	AccessToken TokenKind = "access"
	InviteToken TokenKind = "invite"
)

func NewTokenSigner(kind TokenKind) (*TokenSigner, error) {
	var options TokenOptions
	if err := viper.UnmarshalKey("jwt:"+string(kind), &options); err != nil {
		return nil, err
	}

	return &TokenSigner{&options}, nil
}

func (h *TokenSigner) Sign(claims *Claims) (string, *jwt.Token, error) {
	claims.Audience = jwt.ClaimStrings{h.options.Audience}
	claims.Issuer = h.options.Issuer
	claims.ExpiresAt = jwt.NewNumericDate(claims.IssuedAt.Time.Add(h.options.Lifetime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, _ := base64.StdEncoding.DecodeString(h.options.Secret)
	jwt, err := token.SignedString(secret)
	return jwt, token, err
}
