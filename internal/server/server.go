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
	req := &http.Request{}
	var bufErr error
	for {
		if req.DoneReading {
			break
		}
		curReqData := make([]byte, buffSize)
		br, err := conn.Read(curReqData)
		if err != nil {
			log.Println(err)
			bufErr = err
			// TODO: handle this error somehow (send error response?)
			break
		}
		log.Printf("read %d bytes", br)
		if !req.LineIsRead {
			err = req.ReadLine(&curReqData)
			if err != nil {
				log.Println(err)
				bufErr = err
				// TODO: handle this error somehow (send error response?)
				break
			}
		}
		if !req.HeadersAreRead {
			err = req.ReadHeaders(&curReqData)
			if err != nil {
				log.Println(err)
				bufErr = err
				// TODO: handle this error somehow (send error response?)
				break
			}
		}
		if !req.BodyIsRead {
			err = req.ReadBody(&curReqData, br)
			if err != nil {
				log.Println(err)
				bufErr = err
				// TODO: handle this error somehow (send error response?)
				break
			}
		}
	}

	if bufErr != nil {
		if errors.Is(bufErr, http.ErrInvalidRequestLine) {
			writeErrResponse(conn, http.StatusBadRequest)
		} else {
			writeErrResponse(conn, http.StatusInternalServerError)
		}
		log.Println(bufErr)
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
		writeErrResponse(conn, http.StatusInternalServerError)
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

func writeErrResponse(conn net.Conn, code int) {
	msg, _ := http.GetStatusCodeMessage(code)
	bytes := http.BuildResponseBytes(http.SendErrorResponse(code, msg))
	_, err := conn.Write(bytes)
	if err != nil {
		log.Printf("error writing error response %s", err)
	}
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
