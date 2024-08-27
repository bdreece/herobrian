//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Status -linecomment
package systemd

import (
	"fmt"
	"strings"
)

type Status int

const (
	StatusActiveRunning Status = iota // active (running)
	StatusActiveExited                // active (exited)
	StatusActiveWaiting               // active (waiting)
	StatusInactive                    // inactive
	StatusEnabled                     // enabled
	StatusDisabled                    // disabled
	StatusFailed                      // failed
	StatusStatic                      // static
	StatusMasked                      // masked
	StatusAlias                       // alias
	StatusLinked                      // linked
)

func (s *Status) UnmarshalText(data []byte) error {
	lines := strings.Split(string(data), "\n")
	if len(lines) < 5 {
		return fmt.Errorf("received invalid systemd status output %q", string(data))
	}

	activeLine := lines[4]
	switch true {
	case strings.Contains(activeLine, StatusActiveRunning.String()):
		*s = StatusActiveRunning
	case strings.Contains(activeLine, StatusActiveExited.String()):
		*s = StatusActiveExited
	case strings.Contains(activeLine, StatusActiveWaiting.String()):
		*s = StatusActiveWaiting
	case strings.Contains(activeLine, StatusInactive.String()):
		*s = StatusInactive
	case strings.Contains(activeLine, StatusEnabled.String()):
		*s = StatusEnabled
	case strings.Contains(activeLine, StatusDisabled.String()):
		*s = StatusDisabled
	case strings.Contains(activeLine, StatusStatic.String()):
		*s = StatusStatic
	case strings.Contains(activeLine, StatusMasked.String()):
		*s = StatusMasked
	case strings.Contains(activeLine, StatusLinked.String()):
		*s = StatusLinked
	default:
		return fmt.Errorf("received unknown systemd active value %q", activeLine)
	}

	return nil
}
