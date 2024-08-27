package linode

import (
	"fmt"

	"go.uber.org/config"
)

type Options struct {
	InstanceID  string `yaml:"instance_id"`
	AccessToken string `yaml:"access_token"`
}

func Configure(provider config.Provider) (*Options, error) {
	opts := new(Options)
	if err := provider.Get("linode").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to bind linode options: %w", err)
	}

	return opts, nil
}
