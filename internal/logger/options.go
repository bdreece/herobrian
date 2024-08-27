package logger

import (
	"fmt"
	"log/slog"

	"go.uber.org/config"
)

type Options struct {
	slog.HandlerOptions

	Directory string `yaml:"directory"`
}

func Configure(provider config.Provider) (*Options, error) {
	opts := new(Options)
	if err := provider.Get("log").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to bind logger options: %w", err)
	}

	return opts, nil
}
