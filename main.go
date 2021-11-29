package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var debugF = flag.Bool("d", false, "Enable socket level debugging (if supported)")
var ttlF = flag.Int("f", 1, "Specify with what TTL to start. Defaults to 1")
var hopsF = flag.Int("m", 30, "Specify the maximum number of hops (max time-to-live value) the program will probe. The default is 30")
var portF = flag.Int("p", 34500, "Specify the destination port to use. This number will be incremented by each probe")
var probesF = flag.Int("q", 3, "Sets the number of probe packets per hop. The default number is 3")

func main() {
	log.SetPrefix("rt: ")
	log.SetFlags(0)
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage of rt: [-d -f -m -p -q] host\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(flag.Args()) == 0 {
		log.Printf("A host is required\n")
		flag.Usage()
	}
	if len(flag.Args()) > 1 {
		log.Printf("only 1 destination must be specified\n")
		flag.Usage()
	}

	destination := flag.Args()[0]
	addrs, err := net.LookupHost(destination)
	if err != nil {
		log.Fatalf("lookup for %s failed", destination)
	}
	if len(addrs) == 0 {
		log.Fatalf("no addresses were found for %s", destination)
	}
}
