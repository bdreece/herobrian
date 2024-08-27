package web

import (
	"embed"
	"io/fs"
)

var (
	//go:embed templates/*
	fsys         embed.FS
	Templates, _ = fs.Sub(fsys, "templates")
)
