package host

import "context"

type (
	Info struct {
		ID           string        `json:"-"`
		Architecture string        `json:"arch"`
		DNSName      string        `json:"-"`
		Image        *ImageInfo    `json:"image"`
		IPAddress    string        `json:"-"`
		Platform     string        `json:"platform"`
		Processor    ProcessorInfo `json:"processor"`
	}

	ImageInfo struct {
		Type string `json:"type"`
		// Memory size in MiB.
		Memory *int64 `json:"memory"`
		// Storage space in GiB.
		Storage *int64 `json:"storage"`
		Network string `json:"network"`
	}

	ProcessorInfo struct {
		CoreCount      *int32 `json:"cores"`
		ThreadsPerCore *int32 `json:"threads"`
	}

	Provider interface {
		Hosts(ctx context.Context, ids ...string) ([]*Info, error)
		CheckHosts(ctx context.Context, ids ...string) (map[string]string, error)
		StartHosts(ctx context.Context, ids ...string) error
		StopHosts(ctx context.Context, ids ...string) error
		RestartHosts(ctx context.Context, ids ...string) error
	}
)
