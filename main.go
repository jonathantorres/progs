package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	// TODO: figure out which path to use for the configuration file
	// either from the -conf option, or configured from the build
	// the -conf option would override any location set in the build
	c, err := conf.Load(confF)
	if err != nil {
		log.Fatalf("%s, exiting...", err)
	}
	go sigHandler()
	// start the server
	if err := server.Start(c); err != nil {
		log.Fatalf("%s, exiting...", err)
	}
}

func sigHandler() {
	// TODO: goroutine that will be waiting for any signals
	// that will tell the server to
	// reload the configuration files (HUP)
	// or to gracefully shutdown (TERM, INT, QUIT)
	sigShutdown := make(chan os.Signal, 1)
	signal.Notify(sigShutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sigReload := make(chan os.Signal, 1)
	signal.Notify(sigShutdown, syscall.SIGHUP)

	select {
	case <-sigShutdown:
		// TODO: shutdown the server
	case <-sigReload:
		// TODO: reload configuration files
	}
}
