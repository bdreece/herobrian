package linode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseAddr string = "https://api.linode.com/v4"

type Client struct {
    http.Client

    opts *Options
}

func (c *Client) InstanceStatus(ctx context.Context) (string, error) {
    uri := baseAddr + "/linode/instances/" + c.opts.InstanceID
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
    if err != nil {
        return "", err
    }

    res, err := c.Do(req)
    if err != nil {
        return "", err
    }

    if res.StatusCode != 200 {
        return "", fmt.Errorf("invalid status: %q", res.Status)
    }

    var dto struct {
        Status string `json:"status"`
    }

    defer res.Body.Close()
    if err = json.NewDecoder(res.Body).Decode(&dto); err != nil {
        return "", err
    }

    return dto.Status, nil
}

func (c *Client) BootInstance(ctx context.Context) error {
    uri := baseAddr + "/linode/instances/" + c.opts.InstanceID + "/boot"
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, http.NoBody)
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    res, err := c.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode != 200 {
        return fmt.Errorf("invalid status: %s", res.Status)
    }

    return nil
}

func (c *Client) ShutdownInstance(ctx context.Context) error {
    uri := baseAddr + "/linode/instances/" + c.opts.InstanceID + "/shutdown"
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, http.NoBody)
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    res, err := c.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode != 200 {
        return fmt.Errorf("invalid status: %s", res.Status)
    }

    return nil
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
    r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.opts.AccessToken))
    r.Header.Add("Accept", "application/json")
    return c.Client.Do(r)
}

func NewClient(opts *Options) *Client {
    return &Client{
        opts: opts,
    }
}

