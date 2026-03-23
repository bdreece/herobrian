package host

import (
	"context"

	"github.com/bdreece/herobrian/pkg/minecraft/instance"
	"github.com/spf13/viper"
)

type Client struct {
	*Config

	hostname string
	provider Provider
}

func NewClient(hostname string, provider Provider) (*Client, error) {
	cfg := new(Config)
	if err := viper.UnmarshalKey("minecraft:"+hostname, cfg); err != nil {
		return nil, err
	}

	return &Client{cfg, hostname, provider}, nil
}

func (c *Client) Describe(ctx context.Context) (*Info, error) {
	infos, err := c.provider.DescribeHosts(ctx, c.ID)
	if err != nil {
		return nil, err
	}

	return infos[c.ID], nil
}

func (c *Client) Check(ctx context.Context) (Status, error) {
	statuses, err := c.provider.CheckHosts(ctx, c.ID)
	if err != nil {
		return Status(-1), err
	}

	return statuses[c.ID], nil
}

func (c *Client) Start(ctx context.Context) error {
	return c.provider.StartHosts(ctx, c.ID)
}

func (c *Client) Stop(ctx context.Context) error {
	return c.provider.StopHosts(ctx, c.ID)
}

func (c *Client) Restart(ctx context.Context) error {
	return c.provider.RestartHosts(ctx, c.ID)
}

func (c *Client) Instances(ctx context.Context) (instance.Provider, error) {
	return c.provider.Instances(ctx, c.hostname)
}
