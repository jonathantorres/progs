package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jonathantorres/voy/internal/conf"
	"github.com/jonathantorres/voy/internal/server"
)

var versionFlag = flag.Bool("version", false, "print current version")

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Fprintf(os.Stdout, "voy server v0.1.0\n")
		os.Exit(0)
	}
	if err := conf.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
