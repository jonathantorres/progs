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

const (
	versionFlagDesc = "print current version"
	confFlagDesc    = "specify the location of the configuration file"
)

var versionFlag bool
var confFlag string
var defaultConfFile = "/usr/local/etc/voy"

func init() {
	flag.BoolVar(&versionFlag, "version", false, versionFlagDesc)
	flag.BoolVar(&versionFlag, "v", false, versionFlagDesc+"(shorthand)")
	flag.StringVar(&confFlag, "conf", defaultConfFile, confFlagDesc)
	flag.StringVar(&confFlag, "c", defaultConfFile, confFlagDesc+"(shorthand)")
}

func main() {
	flag.Parse()
	if versionFlag {
		fmt.Fprintf(os.Stdout, "voy server v%s\n", voy.Version)
		os.Exit(0)
	}
	c, err := conf.Load(confFlag)
	if err != nil {
		log.Fatalf("%s, exiting...", err)
	}
	if err := server.Start(c); err != nil {
		log.Fatalf("%s, exiting...", err)
	}
}
