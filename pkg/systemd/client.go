package systemd

import (
	"context"
)

type Client interface {
	Status(context.Context, Unit) (*Status, error)
	Enable(context.Context, Unit) error
	Disable(context.Context, Unit) error
	Start(context.Context, Unit) error
	Stop(context.Context, Unit) error
	Restart(context.Context, Unit) error
}

type ClientOptions[T any] struct {
	Transport T      `yaml:"transport"`
	Units     []Unit `yaml:"units"`
}
