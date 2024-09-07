package systemd

import (
	"context"
	"fmt"
	"sync"

	"github.com/melbahja/goph"
	"go.uber.org/config"
)

type sshClient struct {
	*goph.Client
    mu sync.Mutex
}

type SSH struct {
	User    string `yaml:"user"`
	Address string `yaml:"address"`
	KeyPath string `yaml:"key_path"`
}

// Key loads the SSH key for the configured remote user.
func (transport SSH) Key() (goph.Auth, error) {
	return goph.Key(transport.KeyPath, "")
}

func ConfigureSSH(provider config.Provider) (*ClientOptions[SSH], error) {
	opts := new(ClientOptions[SSH])
	if err := provider.Get("systemd").Populate(opts); err != nil {
		return nil, fmt.Errorf("failed to configure SSH systemd client options: %w", err)
	}

	return opts, nil
}

func (c *sshClient) Status(ctx context.Context, unit Unit) (*Status, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
	out, _ := c.Run(fmt.Sprintf("systemctl status %s", unit))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get systemd service status %q: %w", unit, err)
	// }

	var status Status
	if err := status.UnmarshalText(out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status text %q: %w", string(out), err)
	}

	return &status, nil
}

func (c *sshClient) Enable(ctx context.Context, unit Unit) error {
    c.mu.Lock()
    defer c.mu.Unlock()
	_, err := c.Run(fmt.Sprintf("sudo systemctl enable %s", unit))
	if err != nil {
		return fmt.Errorf("failed to enable systemd service %q: %w", unit, err)
	}

	return nil
}

func (c *sshClient) Disable(ctx context.Context, unit Unit) error {
    c.mu.Lock()
    defer c.mu.Unlock()
	_, err := c.Run(fmt.Sprintf("sudo systemctl disable %s", unit))
	if err != nil {
		return fmt.Errorf("failed to disable systemd service %q: %w", unit, err)
	}

	return nil
}

func (c *sshClient) Start(ctx context.Context, unit Unit) error {
    c.mu.Lock()
    defer c.mu.Unlock()
	_, err := c.Run(fmt.Sprintf("sudo systemctl stop %s", unit))
	if err != nil {
		return fmt.Errorf("failed to start systemd service %q: %w", unit, err)
	}

	return nil
}

func (c *sshClient) Stop(ctx context.Context, unit Unit) error {
    c.mu.Lock()
    defer c.mu.Unlock()
	_, err := c.Run(fmt.Sprintf("sudo systemctl stop %s", unit))
	if err != nil {
		return fmt.Errorf("failed to stop systemd service %q: %w", unit, err)
	}

	return nil
}

func (c *sshClient) Restart(ctx context.Context, unit Unit) error {
    c.mu.Lock()
    defer c.mu.Unlock()
	_, err := c.Run(fmt.Sprintf("sudo systemctl restart %s", unit))
	if err != nil {
		return fmt.Errorf("failed to restart systemd service %q: %w", unit, err)
	}

	return nil
}

func NewSSH(opts *ClientOptions[SSH]) (Client, error) {
	key, err := opts.Transport.Key()
	if err != nil {
		return nil, fmt.Errorf("failed to load ssh key %q: %w", opts.Transport.KeyPath, err)
	}

	c, err := goph.New(opts.Transport.User, opts.Transport.Address, key)
	if err != nil {
		return nil, fmt.Errorf("failed to dial remote address %q: %w", opts.Transport.Address, err)
	}

    return &sshClient{Client: c}, err
}
