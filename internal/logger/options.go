package logger

import (
	"fmt"

	"go.uber.org/config"
)

type Options struct {
	AddSource bool   `yaml:"add_source"`
	Level     int    `yaml:"level"`
	Directory string `yaml:"directory"`
}

func Configure(provider config.Provider) (*Options, error) {
	opts := new(Options)
	if err := provider.Get("logger").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to bind logger options: %w", err)
	}

	return opts, nil
}
