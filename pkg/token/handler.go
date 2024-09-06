package token

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
)

type (
	PasswordResetClaims struct {
		UserID int `mapstructure:"sub"`
	}

	UserInviteClaims struct {
		RoleID int `mapstructure:"role"`
	}

	Handler[C any] interface {
		Sign(*C) (string, error)
		Verify(string) (*C, error)
	}

	handler[C any] struct {
		opts   *Options
		parser *jwt.Parser
	}
)

func (handler *handler[C]) Sign(claims *C) (string, error) {
	id, _ := uuid.NewV4()
	now := time.Now()
	validFor, err := time.ParseDuration(handler.opts.ValidFor)
	if err != nil {
		return "", err
	}

	jwtClaims := jwt.MapClaims{
		"jti": id,
		"aud": handler.opts.Audience,
		"iss": handler.opts.Issuer,
		"exp": jwt.NewNumericDate(now.Add(validFor)),
	}

	if err := mapstructure.Decode(claims, &jwtClaims); err != nil {
		return "", fmt.Errorf("failed to encode claims: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	secret, err := base64.StdEncoding.DecodeString(handler.opts.SecretKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(secret)
}

func (handler *handler[C]) Verify(token string) (*C, error) {
	t, err := handler.parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return base64.StdEncoding.DecodeString(handler.opts.SecretKey)
	})
	if err != nil {
		return nil, err
	}

	jwtClaims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to cast claims to jwt.MapClaims")
	}

	claims := new(C)
	if err := mapstructure.Decode(jwtClaims, claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func NewHandler[C any](opts *Options) Handler[C] {
	return &handler[C]{
		opts: opts,
		parser: jwt.NewParser(
			jwt.WithAudience(opts.Audience),
			jwt.WithExpirationRequired(),
			jwt.WithIssuedAt(),
			jwt.WithIssuer(opts.Issuer),
		),
	}
}
