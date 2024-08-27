package logger

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

func New(opts *Options) (*slog.Logger, error) {
	const (
		flag int         = os.O_WRONLY | os.O_CREATE | os.O_APPEND
		perm fs.FileMode = 0o0644
	)

	path := filepath.Join(opts.Directory, "herobrian.log")
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %q: %w", path, err)
	}

	w := io.MultiWriter(f, os.Stdout)
	handler := slog.NewTextHandler(w, &opts.HandlerOptions)
	return slog.New(handler), nil
}
