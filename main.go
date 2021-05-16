package main

import (
	"encoding/binary"
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
const headerSize = 8

var (
	packetSize     = 56   // the number of  bytes to be sent, the -s flag can change this
	recvBufferSize = 1024 // buffer size when receiving replies
	packetId       = 0    // id for each packet sent
	numTransmitted = 0    // number of packets sent
)

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
	destination := flag.Args()[0]
	addrs, err := net.LookupHost(destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "zing: lookup for %s failed\n", destination)
		os.Exit(1)
	}
	if len(addrs) == 0 {
		fmt.Fprintf(os.Stderr, "zing: no addresses were found for %s\n", destination)
		os.Exit(1)
	}
	solvedDest := addrs[0]
	conn, err := connect(solvedDest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "zing: error connecting: %s\n", err)
		os.Exit(1)
	}

	packetId = os.Getpid() & 0xffff
	printPingMessage(destination, solvedDest)
	go pinger(conn)
	go recvPing(conn)

	// TODO: create signal handler that will terminate
	// the program when a SIGINT is sent to the process (^C)
	// simulate wait for now
	wait := make(chan struct{})
	<-wait
}

type packet struct {
	pType    uint8
	code     uint8
	checksum uint16
	id       uint16
	seqNum   uint16
	data     []byte
}

func newPacket(id uint16, seq uint16) *packet {
	return &packet{
		pType:  uint8(8),
		code:   uint8(0),
		id:     id,
		seqNum: seq,
		data:   nil,
	}
}

func (p *packet) buildData() []byte {
	pData := make([]byte, headerSize+packetSize)
	pData[0], pData[1] = byte(p.pType), byte(p.code)       // type and code
	pData[2], pData[3] = byte(0), byte(0)                  // checksum
	pData[4], pData[5] = byte(p.id>>8), byte(p.id)         // id
	pData[6], pData[7] = byte(p.seqNum>>8), byte(p.seqNum) // seq number

	garbageDataIdx := headerSize
	packSize := packetSize

	// store the timestamp if we can
	if packSize >= 8 {
		b := binary.PutVarint(pData[garbageDataIdx:], time.Now().UnixNano())
		packSize -= b
		garbageDataIdx += b
	}

	// build packet data
	for i := garbageDataIdx; i < packSize; i++ {
		pData[i] = byte(0) // todo: fill with random ascii characters
	}
	p.data = pData[headerSize:]
	csum := calculateChecksum(pData)
	p.checksum = csum
	pData[2] = byte(csum >> 8)
	pData[3] = byte(csum & 255)

	return pData
}

func printPingMessage(destination, solvedDest string) {
	fmt.Fprintf(os.Stdout, "PING %s ", destination)
	if solvedDest != "" {
		fmt.Fprintf(os.Stdout, "(%s)", solvedDest)
	}
	fmt.Fprintf(os.Stdout, " %d bytes of data.\n", packetSize)
}

func printUsage() {
	// TODO
}

func pinger(conn net.Conn) {
	for {
		if err := sendPingPacket(conn); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func connect(dest string) (net.Conn, error) {
	raddr := net.IPAddr{
		IP: net.ParseIP(dest),
	}
	conn, err := net.DialIP("ip4:1", nil, &raddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func sendPingPacket(conn net.Conn) error {
	numTransmitted++
	pack := newPacket(uint16(packetId), uint16(numTransmitted))
	b, err := conn.Write(pack.buildData())
	if err != nil {
		return err
	}
	fmt.Printf("sent %d bytes to %s\n", b, conn.RemoteAddr().String())
	return nil
}

func recvPing(conn net.Conn) {
	// this will receive the reply messages from the echo requests
	buf := make([]byte, recvBufferSize)
	for {
		b, err := conn.Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			continue
		}
		fmt.Printf("received %d bytes from %s\n", b, conn.RemoteAddr().String())
	}
}

// todo: fix this calculation, it only works when bytes with 0 are sent in the payload
func calculateChecksum(msg []byte) uint16 {
	// build out the data in 16-bit chunks
	words := make([]uint16, 0, (headerSize+packetSize)/2)
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
