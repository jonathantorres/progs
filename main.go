package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jonathantorres/voy/internal/conf"
	"github.com/jonathantorres/voy/internal/server"
	"github.com/jonathantorres/voy/internal/voy"
)

var versionFlag = flag.Bool("version", false, "print current version")
var confFlag = flag.String("conf", defaultConfFile, "specify the location of the configuration file")
var defaultConfFile = "/usr/local/etc/voy"

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Fprintf(os.Stdout, "voy server v%s\n", voy.Version)
		os.Exit(0)
	}
	c, err := conf.Load(*confFlag)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Start(c); err != nil {
		log.Fatal(err)
	}
}
