//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Status -linecomment -trimprefix Status
package linode

import "fmt"

type Status int

const (
	StatusRunning           Status = iota // running
	StatusOffline                         // offline
	StatusBooting                         // booting
	StatusRebooting                       // rebooting
	StatusShuttingDown                    // shutting_down
	StatusStopped                         // stopped
	StatusBillingSuspension               // billing_suspension
)

func (s *Status) UnmarshalText(data []byte) error {
	switch string(data) {
	case StatusRunning.String():
		*s = StatusRunning
	case StatusOffline.String():
		*s = StatusOffline
	case StatusBooting.String():
		*s = StatusBooting
	case StatusRebooting.String():
		*s = StatusRebooting
	case StatusShuttingDown.String():
		*s = StatusShuttingDown
	case StatusStopped.String():
		*s = StatusStopped
	case StatusBillingSuspension.String():
		*s = StatusBillingSuspension
	default:
		return fmt.Errorf("invalid linode instance status %q", string(data))
	}

	return nil
}
