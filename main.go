package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jonathantorres/voy/internal/conf"
	"github.com/jonathantorres/voy/internal/server"
)

const version = "0.1.0"

var versionFlag = flag.Bool("version", false, "print current version")

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Fprintf(os.Stdout, "voy server v%s\n", version)
		os.Exit(0)
	}
	if err := conf.Validate(); err != nil {
		log.Fatal(err)
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
