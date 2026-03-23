package ssh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gorcon/rcon"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/bdreece/herobrian/pkg/minecraft/instance"
)

type Config struct {
	Port       int      `mapstructure:"port"`
	User       string   `mapstructure:"user"`
	Identity   string   `mapstructure:"identity"`
	KnownHosts []string `mapstructure:"known_hosts"`
}

type Provider struct {
	Client    *ssh.Client
	Instances map[string]instance.Config
}

var _ instance.Provider = (*Provider)(nil)

func NewProviderFactory() instance.ProviderFactory {
	return instance.ProviderFactoryFunc(func(ctx context.Context, hostname string) (instance.Provider, error) {
		cfg := new(Config)
		if err := viper.UnmarshalKey("minecraft:"+hostname+":ssh", cfg); err != nil {
			return nil, err
		}

		instances := make(map[string]instance.Config)
		if err := viper.UnmarshalKey("minecraft:"+hostname+":instances", &instances); err != nil {
			return nil, err
		}

		knownhosts, err := knownhosts.New(cfg.KnownHosts...)
		if err != nil {
			return nil, err
		}

		pem, err := os.ReadFile(cfg.Identity)
		if err != nil {
			return nil, err
		}

		key, err := ssh.ParseRawPrivateKey(pem)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			return nil, err
		}

		addr := net.JoinHostPort(hostname, fmt.Sprint(cfg.Port))
		client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
			User:            cfg.User,
			HostKeyCallback: knownhosts,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
		})
		if err != nil {
			return nil, err
		}

		return &Provider{client, instances}, nil
	})
}

// DescribeInstances implements [instance.Describer].
func (provider *Provider) DescribeInstance(ctx context.Context, id string) (map[string]string, error) {
	sess, err := provider.Client.NewSession()
	if err != nil {
		return nil, err
	}

	data, err := sess.Output(fmt.Sprintf("sudo systemctl show --no-pager minecraft@%s", id))
	if err != nil {
		return nil, err
	}

	properties := maps.Collect(func(yield func(string, string) bool) {
		for line := range strings.Lines(string(data)) {
			var key, value string

			if _, err := fmt.Sscanf(line, "%s=%s", &key, &value); err != nil {
				panic(err)
			}

			if !yield(key, value) {
				return
			}
		}
	})

	return properties, nil
}

// StartInstance implements [instance.Starter].
func (provider *Provider) StartInstance(ctx context.Context, id string) error {
	sess, err := provider.Client.NewSession()
	if err != nil {
		return err
	}

	return sess.Run(fmt.Sprintf("systemctl start minecraft@%s", id))
}

// StopInstance implements [instance.Stopper].
func (provider *Provider) StopInstance(ctx context.Context, id string) error {
	sess, err := provider.Client.NewSession()
	if err != nil {
		return err
	}

	return sess.Run(fmt.Sprintf("systemctl stop minecraft@%s", id))
}

// RestartInstance implements [instance.Restarter].
func (provider *Provider) RestartInstance(ctx context.Context, id string) error {
	sess, err := provider.Client.NewSession()
	if err != nil {
		return err
	}

	return sess.Run(fmt.Sprintf("systemctl restart minecraft@%s", id))
}

// TraceInstance implements [instance.Tracer].
func (provider *Provider) TraceInstance(ctx context.Context, id string) (<-chan instance.Trace, error) {
	type trace struct {
		RealtimeTimestamp int64  `json:"__REALTIME_TIMESTAMP"`
		Message           string `json:"MESSAGE"`
	}

	traces := make(chan instance.Trace, 1)
	sess, err := provider.Client.NewSession()
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	sess.Stdout = w

	go func(r io.ReadCloser) {
		defer close(traces)
		defer r.Close() //nolint:errcheck
		dec := json.NewDecoder(r)

		for {
			var t trace

			if ctx.Err() != nil {
				return
			}

			if err := dec.Decode(&t); err != nil {
				panic(err)
			}

			traces <- instance.Trace{
				Timestamp: time.UnixMilli(t.RealtimeTimestamp),
				Message:   t.Message,
			}
		}
	}(r)

	return traces, nil
}

func (provider *Provider) DialInstance(ctx context.Context, id string) (instance.Conn, error) {
	instance, ok := provider.Instances[id]
	if !ok {
		return nil, fmt.Errorf("invalid instance: %q", id)
	}

	conn, err := provider.Client.Dial("tcp", net.JoinHostPort("localhost", fmt.Sprint(instance.RCONPort)))
	if err != nil {
		return nil, err
	}

	return rcon.Open(conn, instance.RCONPassword)
}
