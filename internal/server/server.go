package server

import (
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

const buffSize = 1024

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
			for {
				conn, err := l.Accept()
				if err != nil {
					log.Print(err)
					continue
				}
				go handleConn(conn)
			}
			l.Close()
			log.Printf("goroutine done")
		}(p)
	}
	wg.Wait()
	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reqData := make([]byte, buffSize)
	read, err := conn.Read(reqData)
	if err != nil {
		msg, _ := http.GetStatusCodeMessage(http.StatusInternalServerError)
		bytes := http.BuildResponseBytes(http.SendErrorResponse(http.StatusInternalServerError, msg))
		conn.Write(bytes)
		log.Println(err)
		return
	}
	log.Printf("processing request %d bytes read", read)
	req, err := http.NewRequest(reqData)
	if err != nil {
		if errors.Is(err, http.ErrInvalidRequestLine) {
			msg, _ := http.GetStatusCodeMessage(http.StatusBadRequest)
			bytes := http.BuildResponseBytes(http.SendErrorResponse(http.StatusBadRequest, msg))
			conn.Write(bytes)
		} else {
			msg, _ := http.GetStatusCodeMessage(http.StatusInternalServerError)
			bytes := http.BuildResponseBytes(http.SendErrorResponse(http.StatusInternalServerError, msg))
			conn.Write(bytes)
		}
		log.Println(err)
		return
	}
	log.Printf("%s %s HTTP/%d.%d", req.Method, req.Uri, req.HTTPVersionMajor, req.HTTPVersionMinor)

	code, headers, body, err := processRequest(req)
	if err != nil {
		// TODO: Handle any errors to the client here :)
		log.Println(err)
		return
	}

	res := http.NewResponse(code, headers, body)
	written, err := conn.Write(http.BuildResponseBytes(res))
	if err != nil {
		msg, _ := http.GetStatusCodeMessage(http.StatusInternalServerError)
		bytes := http.BuildResponseBytes(http.SendErrorResponse(http.StatusInternalServerError, msg))
		conn.Write(bytes)
		log.Println(err)
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

func getPortsToListen(conf *conf.Conf) ([]int, error) {
	foundPorts := make([]int, 0)
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
	ports := make([]int, 0)
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
