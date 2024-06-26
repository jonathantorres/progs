package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// starts the server process and handles every request sent to it
// handles server start, restart and shutdown

func Start(conf *Conf) error {
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
			for i := 0; i < conf.Workers; i++ {
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
	req := NewRequest(conn)
	err := req.Parse()
	if err != nil {
		// TODO: FIX!
		if errors.Is(err, ErrInvalidRequestLine) {
			writeErrResponse(conn, StatusBadRequest)
		} else {
			writeErrResponse(conn, StatusInternalServerError)
		}
		return
	}

	code, headers, body, err := processRequest(req)
	if err != nil {
		writeErrResponse(conn, StatusInternalServerError)
		return
	}

	res := NewResponse(code, headers, body)
	_, err = conn.Write(BuildResponseBytes(res))
	if err != nil {
		writeErrResponse(conn, StatusInternalServerError)
		return
	}
}

func processRequest(req *Request) (int, map[string]string, []byte, error) {
	headers := make(map[string]string)
	body := make([]byte, 0)
	code := 200 // TODO

	body = append(body, []byte("Hello, world")...) // TODO
	headers["Content-Type"] = "text/html"          // TODO
	return code, headers, body, nil
}

func writeErrResponse(conn net.Conn, code int) {
	msg, _ := GetStatusCodeMessage(code)
	bytes := BuildResponseBytes(SendErrorResponse(code, msg))
	_, err := conn.Write(bytes)
	if err != nil {
		log.Printf("error writing error response %s", err)
	}
}

func getPortsToListen(conf *Conf) ([]int, error) {
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
