package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	icmpHeaderSize    = 8
	ipHeaderSize      = 20
	defaultPacketSize = 56
)

var (
	packetSize     = defaultPacketSize // the number of  bytes to be sent, the -s flag can change this
	recvBufferSize = 1024              // buffer size when receiving replies
	packetId       = 0                 // id for each packet sent
	numTransmitted = 0                 // number of packets sent
	numReceived    = 0                 // number of packets received
)

var countF = flag.Int("c", 0, "Stop after sending -c packets")
var waitF = flag.Int("i", 1, "Wait -i seconds between sending each packet")
var exitF = flag.Bool("o", false, "Exit successfully after receiving one reply packet")
var packetSizeF = flag.Int("s", defaultPacketSize, "Specify the number of data bytes to be sent")
var timeoutF = flag.Int("t", 0, "Timeout, in seconds before zing exits regardless of how many packets have been received")

var transmissionTimes []float64

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

	if *packetSizeF != defaultPacketSize {
		packetSize = *packetSizeF
	}

	transmissionTimes = make([]float64, 0, 15) // arbitrary value
	packetId = os.Getpid() & 0xffff
	printPingMessage(destination, solvedDest)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT)
	go pinger(conn)
	go recvPing(conn, sig)

	if *timeoutF > 0 {
		go timeout(sig)
	}

	<-sig
	printStats(destination)
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

func timeout(sig chan os.Signal) {
	select {
	case <-time.After(time.Duration(*timeoutF) * time.Second):
		sig <- syscall.SIGQUIT
	}
}

func pinger(conn net.Conn) {
	for {
		if err := sendPingPacket(conn); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			break
		}
		time.Sleep(time.Duration(*waitF) * time.Second)
		if *countF > 0 && numReceived >= *countF {
			break
		}
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
	pack := newPacket(uint16(packetId), uint16(numTransmitted))
	_, err := conn.Write(pack.buildData())
	if err != nil {
		return err
	}
	numTransmitted++
	return nil
}

func recvPing(conn net.Conn, sig chan<- os.Signal) {
	// this will receive the reply messages from the echo requests
	buf := make([]byte, recvBufferSize)
	for {
		b, err := conn.Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			continue
		}
		printReceivedPacket(buf, b, conn)
		if (*countF > 0 && numReceived >= *countF) || (*exitF && numReceived >= 1) {
			sig <- syscall.SIGQUIT
			break
		}
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
		fmt.Printf(" time=%s\n", fmt.Sprintf("%.3fms", packTime))
		transmissionTimes = append(transmissionTimes, packTime)
	}
}

func printStats(destination string) {
	fmt.Println()
	fmt.Printf("--- %s ping statistics ---\n", destination)
	fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", numTransmitted, numReceived, calculatePacketLoss())
	min, max, avg, stddev := calculateAverages()
	fmt.Printf("round-trip min/max/avg/stddev = %.3f/%.3f/%.3f/%.3f ms\n", min, max, avg, stddev)
}

func calculatePacketLoss() float64 {
	return float64((numTransmitted - numReceived) * 100 / numTransmitted)
}

func calculateAverages() (float64, float64, float64, float64) {
	var min, max, avg, stddev float64
	if len(transmissionTimes) == 0 {
		return min, max, avg, stddev
	}

	min = transmissionTimes[0]
	max = transmissionTimes[0]
	var sum float64
	for _, t := range transmissionTimes {
		sum += t
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}
	avg = sum / float64(numReceived)

	// calculate standard deviation
	var variance float64
	for _, t := range transmissionTimes {
		diff := t - avg
		diff = diff * diff
		variance += diff
	}
	stddev = math.Sqrt(variance / float64(numReceived))
	return min, max, avg, stddev
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

func calculatePacketTime(buf []byte) (float64, error) {
	tsBytes := buf[28:37]
	n, v := binary.Varint(tsBytes)
	if v <= 0 {
		return 0.0, fmt.Errorf("error decoding the timestamp: %d\n", v)
	}
	now := time.Now().UnixNano()
	ms := now - n
	return float64(ms) / 1000000.00, nil
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
