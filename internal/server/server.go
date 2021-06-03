package server

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/jonathantorres/voy/internal/conf"
	"github.com/jonathantorres/voy/internal/http"
)

// starts the server process and handles every request sent to it
// handles server start, restart and shutdown

func Start(conf *conf.Conf) error {
	ports, err := getPortsToListen(conf)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, p := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("goroutine listening on %d", port)
			var lwg sync.WaitGroup
			for i := 0; i < 5; i++ { // this number should be configurable
				lwg.Add(1)
				go func(l net.Listener) {
					defer lwg.Done()
					for {
						conn, err := l.Accept()
						if err != nil {
							log.Print(err)
							continue
						}
						go handleConn(conn)
					}
				}(l)
			}
			lwg.Wait()
			l.Close()
			log.Printf("goroutine done")
		}(p)
	}
	wg.Wait()
	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	req := http.NewRequest(conn)
	err := req.Parse()
	if err != nil {
		// TODO: FIX!
		if errors.Is(err, http.ErrInvalidRequestLine) {
			writeErrResponse(conn, http.StatusBadRequest)
		} else {
			writeErrResponse(conn, http.StatusInternalServerError)
		}
		return
	}

	log.Printf("%s %s HTTP/%d.%d", req.Method, req.Uri, req.HTTPVersionMajor, req.HTTPVersionMinor)
	code, headers, body, err := processRequest(req)
	if err != nil {
		// TODO: Handle any errors to the client here :)
		log.Printf("processRequest error: %s\n", err)
		return
	}

	res := http.NewResponse(code, headers, body)
	written, err := conn.Write(http.BuildResponseBytes(res))
	if err != nil {
		writeErrResponse(conn, http.StatusInternalServerError)
		log.Printf("conn.Write error: %s\n", err)
	}
	log.Printf("request processed %d bytes written", written)
	log.Printf("HTTP/%d.%d %d %s", res.HTTPVersionMajor, res.HTTPVersionMinor, res.Code, res.Message)
}

func processRequest(req *http.Request) (int, map[string]string, []byte, error) {
	headers := make(map[string]string)
	body := make([]byte, 0)
	code := 200 // TODO

	body = append(body, []byte("Hello, world")...) // TODO
	headers["Content-Type"] = "text/html"          // TODO
	return code, headers, body, nil
}

func writeErrResponse(conn net.Conn, code int) {
	msg, _ := http.GetStatusCodeMessage(code)
	bytes := http.BuildResponseBytes(http.SendErrorResponse(code, msg))
	_, err := conn.Write(bytes)
	if err != nil {
		log.Printf("error writing error response %s", err)
	}
}

func getPortsToListen(conf *conf.Conf) ([]int, error) {
	foundPorts := make([]int, 0, 5)
	if conf.DefaultServer != nil {
		for _, p := range conf.DefaultServer.Ports {
			foundPorts = append(foundPorts, p)
		}
	}
	if conf.Vhosts != nil && len(conf.Vhosts) != 0 {
		for _, vhost := range conf.Vhosts {
			if vhost.Ports != nil {
				for _, p := range vhost.Ports {
					foundPorts = append(foundPorts, p)
				}
			}
		}
	}
	if len(foundPorts) == 0 {
		return nil, errors.New("no ports to listen")
	}
	// don't allow duplicated port numbers
	ports := make([]int, 0, len(foundPorts))
	for _, fp := range foundPorts {
		portFound := false
		if len(ports) > 0 {
			for _, p := range ports {
				if fp == p {
					portFound = true
				}
			}
		}
		if !portFound {
			ports = append(ports, fp)
		}
	}
	return ports, nil
}

func readTheLine(b []byte) bool {
	res := bytes.Split(b, []byte(" "))
	if len(res) == 3 {
		return true
	}
	return false
}
