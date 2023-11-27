package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/ipv6"
)

const (
	icmpHeaderSize    = 8
	ipHeaderSize      = 20
	defaultPacketSize = 56
)

var (
	packetSize     = defaultPacketSize // the number of  bytes to be sent, the -s flag can change this
	recvBufferSize = 1024              // buffer size when receiving replies
	packetID       = 0                 // id for each packet sent
	numTransmitted = 0                 // number of packets sent
	numReceived    = 0                 // number of packets received
)

var countF = flag.Int("c", 0, "Stop after sending -c packets")
var debugF = flag.Bool("d", false, "Set the SO_DEBUG option on the socket being used")
var waitF = flag.Int("i", 1, "Wait -i seconds between sending each packet")
var exitF = flag.Bool("o", false, "Exit successfully after receiving one reply packet")
var ip4F = flag.Bool("4", false, "Use IPv4 only")
var ip6F = flag.Bool("6", false, "Use IPv6 only")
var packetSizeF = flag.Int("s", defaultPacketSize, "Specify the number of data bytes to be sent")
var timeoutF = flag.Int("t", 0, "Timeout, in seconds before ping exits regardless of how many packets have been received")

var transmissionTimes []float64

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of ping:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "ping: a destination must be specified\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(flag.Args()) > 1 {
		fmt.Fprintf(os.Stderr, "ping: only 1 destination must be specified\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	destination := flag.Args()[0]
	addrs, err := net.LookupHost(destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ping: lookup for %s failed\n", destination)
		os.Exit(1)
	}
	if len(addrs) == 0 {
		fmt.Fprintf(os.Stderr, "ping: no addresses were found for %s\n", destination)
		os.Exit(1)
	}
	solvedDest, err := getIPAddr(addrs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ping: error resolving address: %s\n", err)
		os.Exit(1)
	}
	conn, err := connect(solvedDest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ping: error connecting: %s\n", err)
		os.Exit(1)
	}

	if *packetSizeF != defaultPacketSize {
		packetSize = *packetSizeF
	}

	transmissionTimes = make([]float64, 0, 15) // arbitrary value
	packetID = os.Getpid() & 0xffff
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
	ipv6     bool
}

func newPacket(id uint16, seq uint16, ipv6 bool) *packet {
	typ := uint8(8)
	if ipv6 {
		typ = uint8(128)
	}
	return &packet{
		pType:  typ,
		code:   uint8(0),
		id:     id,
		seqNum: seq,
		data:   nil,
		ipv6:   ipv6,
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

func isIPv6(addr string) bool {
	var is6 bool
	for _, a := range addr {
		switch a {
		case ':':
			is6 = true
			break
		}
	}
	return is6
}

func isIPv4(addr string) bool {
	var is4 bool
	for _, a := range addr {
		switch a {
		case '.':
			is4 = true
			break
		}
	}
	return is4
}

func printPingMessage(destination string, solvedDest net.IP) {
	fmt.Fprintf(os.Stdout, "PING %s ", destination)
	if solvedDest.String() != "" {
		fmt.Fprintf(os.Stdout, "(%s)", solvedDest.String())
	}
	fmt.Fprintf(os.Stdout, " %d bytes of data.\n", packetSize)
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

func connect(dest net.IP) (net.Conn, error) {
	ipVer := 4
	protoNum := 1
	if *ip6F {
		ipVer = 6
		protoNum = 58
	}
	raddr := net.IPAddr{
		IP: dest,
	}
	conn, err := net.DialIP(fmt.Sprintf("ip%d:%d", ipVer, protoNum), nil, &raddr)
	if err != nil {
		return nil, err
	}
	if *debugF {
		err = setSocketDebugOption(conn)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func getIPAddr(addrs []string) (net.IP, error) {
	for _, a := range addrs {
		pa := net.ParseIP(a)
		if pa == nil {
			continue // ignore invalid addresses
		}
		if *ip6F {
			// we are only interested in IPv6 addresses
			if isIPv6(a) {
				return pa, nil
			}
		} else if *ip4F {
			// we are only interested in IPv4 addresses
			if isIPv4(a) {
				return pa, nil
			}
		} else {
			// we don't care which IP type,
			// but if this hostname resolves to an IPv6 address,
			// enable the flag so that things work as expected
			if isIPv6(a) {
				*ip6F = true
			}
			return pa, nil
		}
	}
	return nil, fmt.Errorf("address not found")
}

func sendPingPacket(conn net.Conn) error {
	var ipv6 bool
	if *ip6F {
		ipv6 = true
	}
	pack := newPacket(uint16(packetID), uint16(numTransmitted), ipv6)
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
		if err := conn.SetReadDeadline(time.Now().Add(time.Duration((*waitF * 2)) * time.Second)); err != nil {
			fmt.Fprintf(os.Stderr, "deadline error: %s\n", err)
			continue
		}
		b, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Fprintf(os.Stderr, "request timeout: %s\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "read error: %s\n", err)
			}
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
	id := getPacketID(buf)
	// do nothing since this packet does not belong to this process
	if int(id) != packetID {
		return
	}
	typ := getPacketType(buf)
	// make sure we receive only reply packets
	if typ != 129 && typ != 0 {
		return
	}
	numReceived++
	bLen := bytesRead
	if !*ip6F {
		bLen -= ipHeaderSize // ipv4 includes the IP header
	}
	raddr := conn.RemoteAddr().String()
	seq := getPacketSeqNum(buf)
	ttl := int(buf[8])
	if *ip6F {
		ttl = getHopLimit(conn)
	}
	fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d", bLen, raddr, seq, ttl)
	packTime, err := calculatePacketTime(buf)
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

func getPacketID(buf []byte) uint16 {
	st := 24
	en := 26
	if *ip6F {
		st = 4
		en = 6
	}
	packID := buf[st:en]
	id := uint16(packID[0]) << 8
	id |= uint16(packID[1])
	return id & 0xffff
}

func getPacketType(buf []byte) int {
	i := 20
	if *ip6F {
		i = 0
	}
	return int(buf[i])
}

func getPacketSeqNum(buf []byte) uint16 {
	st := 26
	en := 28
	if *ip6F {
		st = 6
		en = 8
	}
	seqNum := buf[st:en]
	num := uint16(seqNum[0]) << 8
	num |= uint16(seqNum[1])
	return num
}

func calculatePacketTime(buf []byte) (float64, error) {
	st := 28
	en := 37
	if *ip6F {
		st = 8
		en = 17
	}
	tsBytes := buf[st:en]
	n, v := binary.Varint(tsBytes)
	if v <= 0 {
		return 0.0, fmt.Errorf("error decoding the timestamp: %d", v)
	}
	now := time.Now().UnixNano()
	ms := now - n
	return float64(ms) / 1000000.00, nil
}

func getHopLimit(conn net.Conn) int {
	c, ok := conn.(net.PacketConn)
	if !ok {
		return 0
	}
	pc := ipv6.NewPacketConn(c)
	hl, err := pc.HopLimit()
	if err != nil {
		return 0
	}
	return hl
}

func setSocketDebugOption(conn *net.IPConn) error {
	rc, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	return rc.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_DEBUG, 1)
	})
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
