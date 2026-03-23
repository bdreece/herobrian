package instance

import "context"

type Client struct {
	id       string
	provider Provider
}

func NewClient(id string, provider Provider) *Client {
	return &Client{id, provider}
}

func (c *Client) Describe(ctx context.Context) (map[string]string, error) {
	return c.provider.DescribeInstance(ctx, c.id)
}

func (c *Client) Start(ctx context.Context) error {
	return c.provider.StartInstance(ctx, c.id)
}

func (c *Client) Stop(ctx context.Context) error {
	return c.provider.StopInstance(ctx, c.id)
}

func (c *Client) Restart(ctx context.Context) error {
	return c.provider.RestartInstance(ctx, c.id)
}

func (c *Client) Trace(ctx context.Context) (<-chan Trace, error) {
	return c.provider.TraceInstance(ctx, c.id)
}

func (c *Client) Dial(ctx context.Context) (Conn, error) {
	return c.provider.DialInstance(ctx, c.id)
}
