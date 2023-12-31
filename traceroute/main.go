package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var debugF = flag.Bool("d", false, "Enable socket level debugging (if supported)")
var ttlF = flag.Int("f", 1, "Specify with what TTL to start. Defaults to 1")
var hopsF = flag.Int("m", 30, "Specify the maximum number of hops (max time-to-live value) the program will probe. The default is 30")
var portF = flag.Int("p", 34500, "Specify the destination port to use. This number will be incremented by each probe")
var probesF = flag.Int("q", 3, "Sets the number of probe packets per hop. The default number is 3")
var probeTimeoutF = flag.Int("w", 5, "Probe timeout. Determines how long to wait for a response to a probe")
var probeIntF = flag.Int("z", 0, "Minimum amount of time to wait between probes (in seconds). The default is 0")

var ip4F = flag.Bool("4", false, "Use IPv4 only")
var ip6F = flag.Bool("6", false, "Use IPv6 only")

const (
	dataBytesLen = 24    // amount of data sent on the UDP packet
	readBufSize  = 1024  // buffer size when reading data from the ICMP packets
	maxPortNum   = 30000 // max port number that we will use in the UDP packet
)

func main() {
	log.SetPrefix("traceroute: ")
	log.SetFlags(0)
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage of traceroute: [-d -f -m -p -q -w -z] host\n")
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
	if *portF < maxPortNum {
		log.Fatalf("port number must be greater than %d", maxPortNum)
	}

	destination := flag.Args()[0]
	addrs, err := net.LookupHost(destination)
	if err != nil {
		log.Fatalf("lookup for %s failed: %s", destination, err)
	}
	if len(addrs) == 0 {
		log.Fatalf("no addresses were found for %s", destination)
	}
	destinationIP, err := getIPAddr(addrs)
	if err != nil {
		log.Fatalf("IP address not found: %s", err)
	}
	printStart(destination, destinationIP)
	go listenICMP()
	startTrace(destinationIP)
}

type tracePacket struct {
	seqNum int64
	ttl    int64
	ts     int64
}

type probeInfo struct {
	routerIP   net.IP
	routerName string
	icmpType   int
	icmpCode   int
	udpPort    int
}

var probChan chan *probeInfo

func listenICMP() {
	ipVer := 4
	protoNum := 1
	if *ip6F {
		ipVer = 6
		protoNum = 58
	}
	laddr := net.IPAddr{
		IP: nil,
	}
	conn, err := net.ListenIP(fmt.Sprintf("ip%d:%d", ipVer, protoNum), &laddr)
	if err != nil {
		log.Fatalf("error listening for ICPMP packets: %s", err)
	}
	probChan = make(chan *probeInfo)
	for {
		buf := make([]byte, readBufSize)
		var raddr net.Addr
		if *ip6F {
			_, raddr, err = conn.ReadFrom(buf)
		} else {
			_, err = conn.Read(buf)
		}
		if err != nil {
			log.Printf("error reading data: %s", err)
			continue
		}
		pInfo := newProbeInfo(raddr, buf)
		probChan <- pInfo
	}
}

func startTrace(destIP net.IP) {
	port := *portF
	var seqNum int
	var done bool
	for ttl := *ttlF; ttl <= *hopsF; ttl++ {
		if done {
			break
		}
		fmt.Printf("%2d ", ttl)
		var prevRouterName string
		for pro := 0; pro < *probesF; pro++ {
			udpConn, err := connectUDP(destIP, port, ttl)
			if err != nil {
				log.Printf("error connecting: %s", err)
				continue
			}
			if *debugF {
				setSocketDebugOption(udpConn) // ignoring any errors
			}
			seqNum++
			port++
			d := tracePacket{
				seqNum: int64(seqNum),
				ttl:    int64(ttl),
				ts:     time.Now().UnixNano(),
			}
			startTS := d.ts
			_, err = udpConn.Write(getTracePacketData(&d))
			if err != nil {
				log.Printf("error sending data: %s", err)
				continue
			}
			timer := time.NewTimer(time.Duration(*probeTimeoutF) * time.Second)
			var pInfo *probeInfo
			select {
			case pInfo = <-probChan:
				timer.Stop()
			case <-timer.C:
				fmt.Printf("* ")
				continue // continue to the next probe
			}
			// make sure the packet is destined for this process
			if pInfo.udpPort != port-1 {
				continue
			}
			endTS := time.Now().UnixNano()
			if pro == 0 {
				printRouterIP(pInfo)
				prevRouterName = pInfo.routerName
			} else if pInfo.routerName != "" && prevRouterName != pInfo.routerName {
				fmt.Printf("\n   ")
				printRouterIP(pInfo)
				prevRouterName = pInfo.routerName
			}
			fmt.Printf("%.3f ms   ", float64(endTS-startTS)/1000000.00)
			if isPortUnreachable(pInfo) {
				done = true
			}
			// wait interval before sending the next probe
			if *probeIntF > 0 {
				time.Sleep(time.Duration(*probeIntF) * time.Second)
			}
		}
		fmt.Println()
	}
}

func connectUDP(destIP net.IP, port int, ttl int) (*net.UDPConn, error) {
	raddr := net.UDPAddr{
		IP:   destIP,
		Port: port,
	}
	ipVer := 4
	if *ip6F {
		ipVer = 6
	}
	udpConn, err := net.DialUDP(fmt.Sprintf("udp%d", ipVer), nil, &raddr)
	if err != nil {
		return nil, err
	}
	if *ip6F {
		nconn := ipv6.NewConn(udpConn)
		err = nconn.SetHopLimit(ttl)
	} else {
		nconn := ipv4.NewConn(udpConn)
		err = nconn.SetTTL(ttl)
	}
	if err != nil {
		return nil, err
	}
	return udpConn, nil
}

func newProbeInfo(raddr net.Addr, buf []byte) *probeInfo {
	var routerName string
	var icmpType int
	var icmpCode int
	var routerIP net.IP
	if *ip6F {
		icmpType = int(buf[0])
		icmpCode = int(buf[1])
		routerIP = net.ParseIP(raddr.String())
	} else {
		routerIP = net.IPv4(buf[12], buf[13], buf[14], buf[15])
		icmpType = int(buf[20])
		icmpCode = int(buf[21])
	}
	udpPortSli := buf[50:52]
	udpPort := uint16(udpPortSli[0]) << 8
	udpPort |= uint16(udpPortSli[1])
	udpPort &= 0xffff

	names, _ := net.LookupAddr(routerIP.String())
	if len(names) > 0 {
		routerName = names[0]
	}
	return &probeInfo{
		routerIP:   routerIP,
		routerName: routerName,
		icmpType:   icmpType,
		icmpCode:   icmpCode,
		udpPort:    int(udpPort),
	}
}

func printRouterIP(pInfo *probeInfo) {
	routerAddr := pInfo.routerIP.String()
	if pInfo.routerName != "" {
		fmt.Printf("%s", pInfo.routerName)
	} else {
		fmt.Printf("%s", routerAddr)
	}
	fmt.Printf(" (%s)", routerAddr)
	fmt.Printf("  ")
}

func isPortUnreachable(pInfo *probeInfo) bool {
	if *ip6F {
		if pInfo.icmpType == 1 && pInfo.icmpCode == 4 {
			return true
		}
	} else {
		if pInfo.icmpType == 3 && pInfo.icmpCode == 3 {
			return true
		}
	}
	return false
}

func getTracePacketData(data *tracePacket) []byte {
	d := make([]byte, dataBytesLen)
	var n int
	n = binary.PutVarint(d, data.seqNum)
	n = binary.PutVarint(d[n:], data.ttl)
	binary.PutVarint(d[n*2:], data.ts)
	return d
}

func printStart(destination string, destinationIP net.IP) {
	fmt.Printf("traceroute to %s", destination)
	if destinationIP != nil {
		fmt.Printf(" (%s),", destinationIP.String())
	}
	fmt.Printf(" %d hops max, %d byte packets\n", *hopsF, dataBytesLen)
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

func setSocketDebugOption(conn *net.UDPConn) error {
	rc, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	return rc.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_DEBUG, 1)
	})
}
