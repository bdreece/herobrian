package security

import (
	"github.com/go-crypt/crypt"
	"github.com/go-crypt/crypt/algorithm"
	"github.com/go-crypt/crypt/algorithm/argon2"
)

type PasswordHasher interface {
	algorithm.Hash
	algorithm.Decoder
}

type Argon2Hasher struct {
	*argon2.Hasher
	*crypt.Decoder
}

func NewArgon2() (*Argon2Hasher, error) {
	decoder := crypt.NewDecoder()
	if err := argon2.RegisterDecoderArgon2id(decoder); err != nil {
		return nil, err
	}

	hasher, err := argon2.New(argon2.WithProfileRFC9106LowMemory())
	if err != nil {
		return nil, err
	}

	return &Argon2Hasher{hasher, decoder}, nil
}
