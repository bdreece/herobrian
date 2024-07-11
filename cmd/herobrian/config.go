package main

import (
	"fmt"
	"os"

	"go.uber.org/config"
)

func loadConfig(path string) (config.Provider, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file %q: %v", path, err)
    }

    return config.NewYAML(config.Source(f))
}
