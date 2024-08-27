package router

import (
	"fmt"

	"go.uber.org/config"
)

type Options struct {
	StaticDirectory string `yaml:"static_dir"`
	AppDirectory    string `yaml:"app_dir"`
}

func Configure(provider config.Provider) (*Options, error) {
	opts := new(Options)
	if err := provider.Get("router").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure router options: %w", err)
	}

	return opts, nil
}
