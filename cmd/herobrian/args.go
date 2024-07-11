package main

import (
	"flag"
	"fmt"
)

type Args struct {
	Port       int
	ConfigPath string
}

func (a Args) Addr() string { return fmt.Sprintf(":%d", a.Port) }

var args Args

func init() {
    flag.IntVar(&args.Port, "p", 3000, "port")
	flag.StringVar(&args.ConfigPath, "c", "configs/settings.yml", "config path")
	flag.Parse()
}
