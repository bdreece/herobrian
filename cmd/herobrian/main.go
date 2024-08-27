package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bdreece/herobrian"
)

func main() {
	var args herobrian.Args
	defer quit()

	flag.IntVar(&args.Port, "p", 3000, "port")
	flag.StringVar(&args.ConfigPath, "c", "configs/settings.yml", "config path")
	flag.Parse()

	herobrian.New(args).Run()
}

func quit() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "unexpected panic occurred: %v", r)
		os.Exit(1)
	}
}
