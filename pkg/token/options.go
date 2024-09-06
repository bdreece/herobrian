package token

import (
	"fmt"

	"go.uber.org/config"
)

type Options struct {
	Audience  string `yaml:"audience"`
	Issuer    string `yaml:"issuer"`
	ValidFor  string `yaml:"valid_for"`
	SecretKey string `yaml:"secret_key"`
}

func Configure(name string, provider config.Provider) (*Options, error) {
	opts := new(Options)
	if err := provider.Get("token." + name).Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure token options %q: %w", name, err)
	}

	return opts, nil
}
