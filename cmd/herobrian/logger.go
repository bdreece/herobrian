package main

import (
	"log/slog"
	"os"

	"go.uber.org/config"
)

func createLogger(provider config.Provider) (*slog.Logger, error) {
    opts := new(slog.HandlerOptions)
    if err := provider.Get("log").Populate(opts); err != nil {
        return nil, err
    }

    return slog.New(slog.NewTextHandler(os.Stdout, opts)), nil
}
