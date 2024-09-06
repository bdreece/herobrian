package identity

import (
	"encoding/base64"
	"fmt"

	"github.com/gorilla/sessions"
	"go.uber.org/config"
)

type SessionOptions struct {
	SigningKey    string `yaml:"signing_key"`
	EncryptionKey string `yaml:"encryption_key"`
}

func ConfigureSession(provider config.Provider) (*SessionOptions, error) {
	opts := new(SessionOptions)
	if err := provider.Get("session").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure session options: %w", err)
	}

	return opts, nil
}

func NewSessionStore(opts *SessionOptions) (*sessions.CookieStore, error) {
	signingKey, err := base64.StdEncoding.DecodeString(opts.SigningKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signing key: %w", err)
	}

	encryptionKey, err := base64.StdEncoding.DecodeString(opts.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	return sessions.NewCookieStore(signingKey[:32], encryptionKey[:32]), nil
}
