package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims interface {
	jwt.Claims

	SetExpirationTime(exp *jwt.NumericDate)
	SetIssuedAt(iat *jwt.NumericDate)
	SetNotBefore(nbf *jwt.NumericDate)
	SetIssuer(iss string)
	SetSubject(sub string)
	SetAudience(aud jwt.ClaimStrings)
}

type registeredClaims struct {
	jwt.RegisteredClaims
}

// SetAudience implements Claims.
func (r *registeredClaims) SetAudience(aud jwt.ClaimStrings) {
	r.Audience = aud
}

// SetExpirationTime implements Claims.
func (r *registeredClaims) SetExpirationTime(exp *jwt.NumericDate) {
	r.ExpiresAt = exp
}

// SetIssuedAt implements Claims.
func (r *registeredClaims) SetIssuedAt(iat *jwt.NumericDate) {
	r.IssuedAt = iat
}

// SetIssuer implements Claims.
func (r *registeredClaims) SetIssuer(iss string) {
	r.Issuer = iss
}

// SetNotBefore implements Claims.
func (r *registeredClaims) SetNotBefore(nbf *jwt.NumericDate) {
	r.NotBefore = nbf
}

// SetSubject implements Claims.
func (r *registeredClaims) SetSubject(sub string) {
	r.Subject = sub
}

var ErrInvalidClaims = errors.New("token: invalid claims type")
