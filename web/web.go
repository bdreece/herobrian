package web

import (
	"embed"
)

var (
	//go:embed templates/*
	Templates embed.FS
)
