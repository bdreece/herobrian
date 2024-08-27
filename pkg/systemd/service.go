package systemd

import (
	"context"
	"fmt"
)

type Unit struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Instance    string `yaml:"instance"`
}

type Service struct {
	client Client
	unit   Unit
}

func (u Unit) String() string { return fmt.Sprintf("%s@%s.service", u.Name, u.Instance) }

func (svc Service) Unit() Unit { return svc.unit }

func (svc Service) Status(ctx context.Context) (*Status, error) {
	return svc.client.Status(ctx, svc.unit)
}

func (svc Service) Enable(ctx context.Context) error {
	return svc.client.Enable(ctx, svc.unit)
}

func (svc Service) Disable(ctx context.Context) error {
	return svc.client.Disable(ctx, svc.unit)
}

func (svc Service) Start(ctx context.Context) error {
	return svc.client.Start(ctx, svc.unit)
}

func (svc Service) Stop(ctx context.Context) error {
	return svc.client.Stop(ctx, svc.unit)
}

func (svc Service) Restart(ctx context.Context) error {
	return svc.client.Restart(ctx, svc.unit)
}
