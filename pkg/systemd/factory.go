package systemd

import "fmt"

type ServiceFactory struct {
	opts *ClientOptions[SSH]
}

func (f ServiceFactory) Units() []Unit { return f.opts.Units }

func (f *ServiceFactory) Create(instance string) (*Service, error) {
	client, err := NewSSH(f.opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH systemd client: %w", err)
	}

	for _, unit := range f.opts.Units {
		if unit.Instance == instance {
			return &Service{
				client: client,
				unit:   unit,
			}, nil
		}
	}

    return nil, fmt.Errorf("systemd instance not found")
}

func NewServiceFactory(opts *ClientOptions[SSH]) *ServiceFactory {
    return &ServiceFactory{opts}
}
