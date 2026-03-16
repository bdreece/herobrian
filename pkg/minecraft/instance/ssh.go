package instance

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHProvider struct {
	Client *ssh.Client
}

// Instances implements [InstanceProvider].
func (provider *SSHProvider) Instances(ctx context.Context, ids ...string) ([]Info, error) {
	sess, err := provider.Client.NewSession()
	if err != nil {
		return nil, err
	}

	_, _ = sess.Output(fmt.Sprintf("sudo systemctl show %s", strings.Join(ids, " ")))
	panic("unimplemented")
}

// StartInstance implements [InstanceProvider].
func (s *SSHProvider) StartInstance(ctx context.Context, id string) error {
	panic("unimplemented")
}

// StopInstance implements [InstanceProvider].
func (s *SSHProvider) StopInstance(ctx context.Context, id string) error {
	panic("unimplemented")
}

// RestartInstance implements [InstanceProvider].
func (s *SSHProvider) RestartInstance(ctx context.Context, id string) error {
	panic("unimplemented")
}

// TraceInstance implements [InstanceProvider].
func (s *SSHProvider) TraceInstance(ctx context.Context, id string) (*Tracer, error) {
	panic("unimplemented")
}

// DialInstance implements [InstanceProvider].
func (s *SSHProvider) DialInstance(ctx context.Context, id string) (InstanceConn, error) {
	panic("unimplemented")
}

var _ InstanceProvider = (*SSHProvider)(nil)
