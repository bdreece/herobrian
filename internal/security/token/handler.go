package token

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type HandlerOptions struct {
	Audience []string      `mapstructure:"aud"`
	Issuer   string        `mapstructure:"iss"`
	Secret   string        `mapstructure:"secret"`
	Lifetime time.Duration `mapstructure:"lifetime"`
}

type Encoder[C Claims] interface {
	Encode(claims C) (string, *jwt.Token, error)
}

type Decoder[C Claims] interface {
	Decode(jwt string) (C, *jwt.Token, error)
}

type Handler[C Claims] interface {
	Encoder[C]
	Decoder[C]
}

type encoder[C Claims] struct {
	options *HandlerOptions
}

func (e *encoder[C]) Encode(claims C) (string, *jwt.Token, error) {
	iat, err := claims.GetIssuedAt()
	if err != nil {
		return "", nil, err
	}

	claims.SetAudience(jwt.ClaimStrings(e.options.Audience))
	claims.SetIssuer(e.options.Issuer)
	claims.SetExpirationTime(jwt.NewNumericDate(iat.Add(e.options.Lifetime)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, _ := base64.StdEncoding.DecodeString(e.options.Secret)
	jwt, err := token.SignedString(secret)

	return jwt, token, err
}

func (e *encoder[C]) keyFunc(*jwt.Token) (any, error) {
	return []byte(e.options.Secret), nil
}

var _ Encoder[Claims] = (*encoder[Claims])(nil)
