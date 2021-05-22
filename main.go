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
	versionFDesc    = "print current version"
	confFDesc       = "specify the location of the configuration file"
	logFDesc        = "specify the location of the log file"
)

func main() {
	var (
		versionF bool
		confF    string
		logF     string
	)
	flag.BoolVar(&versionF, "version", false, versionFDesc)
	flag.BoolVar(&versionF, "v", false, versionFDesc+"(shorthand)")
	flag.StringVar(&confF, "conf", defaultConfFile, confFDesc)
	flag.StringVar(&confF, "c", defaultConfFile, confFDesc+"(shorthand)")
	flag.StringVar(&logF, "log", defaultLogFile, logFDesc)
	flag.StringVar(&logF, "l", defaultLogFile, logFDesc+"(shorthand)")
	flag.Parse()

	if versionF {
		fmt.Fprintf(os.Stdout, "voy server v%s\n", voy.Version)
		os.Exit(0)
	}
	// TODO: initialize logging mechanism
	c, err := conf.Load(confF)
	if err != nil {
		log.Fatalf("%s, exiting...", err)
	}
	if err := server.Start(c); err != nil {
		log.Fatalf("%s, exiting...", err)
	}
}
