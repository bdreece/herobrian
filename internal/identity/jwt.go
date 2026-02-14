package identity

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-viper/mapstructure/v2"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"
)

type JWTOptions struct {
	Audience string `mapstructure:"aud"`
	Issuer   string `mapstructure:"iss"`
	Secret   string `mapstructure:"secret"`
	Lifetime string `mapstructure:"lifetime"`
}

func (opts *JWTOptions) BindConfig(v *viper.Viper) error {
	return v.UnmarshalKey("jwt", opts)
}

type Claims struct {
	Subject     string `mapstructure:"sub"`
	Audience    string `mapstructure:"aud"`
	Issuer      string `mapstructure:"iss"`
	Expiration  int64  `mapstructure:"exp"`
	IssuedAt    int64  `mapstructure:"iat"`
	NotBefore   int64  `mapstructure:"nbf"`
	ID          string `mapstructure:"jti"`
	FirstName   string `mapstructure:"given_name"`
	LastName    string `mapstructure:"family_name"`
	DisplayName string `mapstructure:"preferred_username"`
	Picture     string `mapstructure:"picture"`
	Role        string `mapstructure:"role"`
}

type TokenSigner struct {
	jwtauth *jwtauth.JWTAuth
}

func NewTokenSigner(jwtauth *jwtauth.JWTAuth) *TokenSigner {
	return &TokenSigner{jwtauth}
}

func (s *TokenSigner) Sign(claims *Claims) (jwt.Token, string, error) {
	var mapclaims map[string]any

	if err := mapstructure.Decode(claims, &mapclaims); err != nil {
		return nil, "", err
	}

	return s.jwtauth.Encode(mapclaims)
}

func FromContext(ctx context.Context) (jwt.Token, *Claims, error) {
	token, mapclaims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return token, nil, err
	}

	claims := new(Claims)
	err = mapstructure.Decode(mapclaims, claims)

	return token, claims, err
}
