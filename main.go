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
	defaultPrefix   = "/usr/local/voy"
	defaultConfFile = defaultPrefix + "/conf/voy.conf"
	defaultLogFile  = defaultPrefix + "/log/voy.log"
	versionFlagDesc = "print current version"
	confFlagDesc    = "specify the location of the configuration file"
	logFlagDesc     = "specify the location of the log file"
)

var versionFlag bool
var confFlag string
var logFlag string

func init() {
	flag.BoolVar(&versionFlag, "version", false, versionFlagDesc)
	flag.BoolVar(&versionFlag, "v", false, versionFlagDesc+"(shorthand)")
	flag.StringVar(&confFlag, "conf", defaultConfFile, confFlagDesc)
	flag.StringVar(&confFlag, "c", defaultConfFile, confFlagDesc+"(shorthand)")
	flag.StringVar(&logFlag, "log", defaultLogFile, logFlagDesc)
	flag.StringVar(&logFlag, "l", defaultLogFile, logFlagDesc+"(shorthand)")
}

func main() {
	flag.Parse()
	if versionFlag {
		fmt.Fprintf(os.Stdout, "voy server v%s\n", voy.Version)
		os.Exit(0)
	}
	// TODO: initialize logging mechanism
	c, err := conf.Load(confFlag)
	if err != nil {
		log.Fatalf("%s, exiting...", err)
	}
	if err := server.Start(c); err != nil {
		log.Fatalf("%s, exiting...", err)
	}
}
