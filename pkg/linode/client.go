package linode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var baseUri, _ = url.Parse("https://api.linode.com")

type Client interface {
    InstanceStatus(ctx context.Context) (*Status, error)
    BootInstance(ctx context.Context) error
    ShutdownInstance(ctx context.Context) error
    RebootInstance(ctx context.Context) error
}

type httpClient struct {
    http.Client

    opts *Options
}

func (client *httpClient) InstanceStatus(ctx context.Context) (*Status, error) {
    uri := *baseUri
    uri.Path = fmt.Sprintf("/v4/linode/instances/%s", client.opts.InstanceID)

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
    if err != nil {
        return nil, err
    }

    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("invalid status: %q", res.Status)
    }

    var dto struct {
        Status string `json:"status"`
    }

    defer res.Body.Close()
    if err = json.NewDecoder(res.Body).Decode(&dto); err != nil {
        return nil, err
    }

    var status Status
    if err = status.UnmarshalText([]byte(dto.Status)); err != nil {
        return nil, err
    }

    return &status, nil
}

func (client *httpClient) BootInstance(ctx context.Context) error {
    uri := *baseUri
    uri.Path = fmt.Sprintf("/v4/linode/instances/%s/boot", client.opts.InstanceID)

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), http.NoBody)
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")
    res, err := client.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode != 200 {
        return fmt.Errorf("invalid status: %s", res.Status)
    }

    return nil
}

func (client *httpClient) ShutdownInstance(ctx context.Context) error {
    uri := *baseUri
    uri.Path = fmt.Sprintf("/v4/linode/instances/%s/shutdown", client.opts.InstanceID)

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), http.NoBody)
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")
    res, err := client.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode != 200 {
        return fmt.Errorf("invalid status: %s", res.Status)
    }

    return nil
}

func (client *httpClient) RebootInstance(ctx context.Context) error {
    uri := *baseUri
    uri.Path = fmt.Sprintf("/linode/instances/%s/reboot", client.opts.InstanceID)

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), http.NoBody)
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")
    res, err := client.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode != 200 {
        return fmt.Errorf("invalid status: %s", res.Status)
    }

    return nil
}

func (client *httpClient) Do(r *http.Request) (*http.Response, error) {
    r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.opts.AccessToken))
    r.Header.Add("Accept", "application/json")
    return client.Client.Do(r)
}

func NewHTTP(opts *Options) Client {
    return &httpClient{
        opts: opts,
    }
}

