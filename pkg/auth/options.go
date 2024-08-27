package auth

import (
	"fmt"

	"go.uber.org/config"
)

type Options struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Configure(provider config.Provider) (*Options, error) {
	opts := new(Options)
	if err := provider.Get("auth").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to bind auth options: %w", err)
	}

	return opts, nil
}
