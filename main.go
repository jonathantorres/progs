package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// Send a packet (ICMP echo request) every second, and wait for a reply
// The ping program contains two logical portions: one transmits an
// ICMP echo request message every second and the other receives
// any echo reply messages that are returned
var destination string
var solvedDest string

var conn *net.IPConn

// the number of data bytes to be sent, the -s flag can change this
var packetSize = 56

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "zing: a destination must be specified\n")
		printUsage()
		os.Exit(1)
	}
	if len(flag.Args()) > 1 {
		fmt.Fprintf(os.Stderr, "zing: only 1 destination must be specified\n")
		printUsage()
		os.Exit(1)
	}
	destination = flag.Args()[0]
	addrs, err := net.LookupHost(destination)
	if err != nil || len(addrs) == 0 {
		fmt.Fprintf(os.Stdout, "zing: lookup for %s failed\n", destination)
		printUsage()
		os.Exit(1)
	}
	solvedDest = addrs[0]
	wait := make(chan struct{})
	printPingMessage()
	go pinger()
	go recvPing()

	// TODO: create signal handler that will terminate
	// the program when a SIGINT is sent to the process (^C)
	// simulate wait for now
	<-wait
}

func printPingMessage() {
	fmt.Fprintf(os.Stdout, "PING %s ", destination)
	if solvedDest != "" {
		fmt.Fprintf(os.Stdout, "(%s)", solvedDest)
	}
	fmt.Fprintf(os.Stdout, " %d bytes of data.\n", packetSize)
}

func printUsage() {
	// TODO
}

func pinger() {
	for {
		if err := sendPingPacket(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func sendPingPacket() error {
	raddr := net.IPAddr{
		IP: net.ParseIP(solvedDest),
	}
	var err error
	conn, err = net.DialIP("ip4:1", nil, &raddr)
	if err != nil {
		return err
	}
	msg := make([]byte, 8+packetSize)
	msg[0], msg[1] = byte(8), byte(0) // type and code
	msg[2], msg[3] = byte(0), byte(0) // checksum
	msg[4], msg[5] = byte(0), byte(0) // id
	msg[6], msg[7] = byte(0), byte(0) // seq number

	// build packet data
	for i, offset := 1, 8; i <= packetSize; i, offset = i+1, offset+1 {
		msg[offset] = byte(0)
	}
	csum := calculateChecksum(msg)
	msg[2] = byte(csum >> 8)
	msg[3] = byte(csum & 255)
	_, err = conn.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

func recvPing() {
	// this will receive the reply messages from the echo requests
	buf := make([]byte, 1024)
	for {
		if conn == nil {
			continue
		}
		b, err := conn.Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		fmt.Printf("we got: %d bytes\n", b)
	}
}

// todo: fix this calculation, it only works when bytes with 0 are sent in the payload
func calculateChecksum(msg []byte) uint16 {
	// build out the data in 16-bit chunks
	words := make([]uint16, 0, packetSize/2)
	for i := 0; i < len(msg); i += 2 {
		l := uint16(msg[i]) << 8
		word := uint16(l) | uint16(msg[i+1])
		words = append(words, word)
	}
	// calculate checksum
	var csum uint16
	for _, w := range words {
		csum = csum + w
	}
	return ^csum
}
