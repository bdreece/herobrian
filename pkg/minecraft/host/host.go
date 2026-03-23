package host

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
)

//go:generate go tool stringer -type=Status -trimprefix=Status
type Status int

const (
	StatusPending Status = 16 * iota
	StatusRunning
	StatusShuttingDown
	StatusTerminated
	StatusStopping
	StatusStopped
)
