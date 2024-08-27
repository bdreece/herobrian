package controller

import (
	"fmt"
	"io"
	"strings"
)

type event struct {
	Event string
	Data  string
}

func (e event) WriteTo(w io.Writer) (int64, error) {
	var total int64

	if e.Event != "" {
		n, err := fmt.Fprintf(w, "event: %s\n", e.Event)
		if err != nil {
			return total, err
		}

		total += int64(n)
	}

	for _, line := range strings.Split(e.Data, "\n") {
		n, err := fmt.Fprintf(w, "data: %s\n", line)
		if err != nil {
			return total, err
		}

		total += int64(n)
	}

	n, err := fmt.Fprint(w, "\n")
	if err != nil {
		return total, err
	}

	return total + int64(n), nil
}
