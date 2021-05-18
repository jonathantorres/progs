package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

// Send a packet (ICMP echo request) every second, and wait for a reply
// The ping program contains two logical portions: one transmits an
// ICMP echo request message every second and the other receives
// any echo reply messages that are returned
const (
	icmpHeaderSize = 8
	ipHeaderSize   = 20
)

var (
	packetSize     = 56   // the number of  bytes to be sent, the -s flag can change this
	recvBufferSize = 1024 // buffer size when receiving replies
	packetId       = 0    // id for each packet sent
	numTransmitted = 0    // number of packets sent
	numReceived    = 0    // number of packets received
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
	pData := make([]byte, icmpHeaderSize+packetSize)
	pData[0], pData[1] = byte(p.pType), byte(p.code)       // type and code
	pData[2], pData[3] = byte(0), byte(0)                  // checksum
	pData[4], pData[5] = byte(p.id>>8), byte(p.id)         // id
	pData[6], pData[7] = byte(p.seqNum>>8), byte(p.seqNum) // seq number

	garbageDataIdx := icmpHeaderSize
	packSize := packetSize

	// store the timestamp if we can
	if packSize >= 8 {
		b := binary.PutVarint(pData[garbageDataIdx:], time.Now().UnixNano())
		packSize -= b
		garbageDataIdx += b
	}

	// build packet data
	rand.Seed(time.Now().UnixNano())
	for i := garbageDataIdx; i < packSize; i++ {
		pData[i] = byte(rand.Intn(127))
	}
	p.data = pData[icmpHeaderSize:]
	csum := calculateChecksum(pData)
	p.checksum = csum
	pData[2], pData[3] = byte(csum&255), byte(csum>>8)

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
	_, err := conn.Write(pack.buildData())
	if err != nil {
		return err
	}
	// fmt.Printf("sent %d bytes to %s\n", b, conn.RemoteAddr().String())
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
		printReceivedPacket(buf, b, conn)
	}
}

func printReceivedPacket(buf []byte, bytesRead int, conn net.Conn) {
	id := getPacketId(buf)
	// do nothing since this packet does not belong to this process
	if int(id) != packetId {
		return
	}
	numReceived++
	bLen := bytesRead - ipHeaderSize
	raddr := conn.RemoteAddr().String()
	seq := getPacketSeqNum(buf)
	packTime, err := calculatePacketTime(buf)
	fmt.Printf("%d bytes from %s: icmp_seq=%d", bLen, raddr, seq)
	if err == nil {
		fmt.Printf(" time=%s\n", packTime)
	}
}

func getPacketId(buf []byte) uint16 {
	packId := buf[24:26]
	id := uint16(packId[0]) << 8
	id |= uint16(packId[1])
	return id & 0xffff
}

func getPacketSeqNum(buf []byte) uint16 {
	seqNum := buf[26:28]
	num := uint16(seqNum[0]) << 8
	num |= uint16(seqNum[1])
	return num
}

func calculatePacketTime(buf []byte) (string, error) {
	tsBytes := buf[28:37]
	n, v := binary.Varint(tsBytes)
	if v <= 0 {
		return "", fmt.Errorf("error decoding the timestamp: %d\n", v)
	}
	now := time.Now().UnixNano()
	ms := now - n
	return fmt.Sprintf("%.3fms", float64(ms)/1000000.00), nil
}

func calculateChecksum(b []byte) uint16 {
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	return ^uint16(s)
}
