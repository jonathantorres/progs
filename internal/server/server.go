package server

import (
	"fmt"
	"log"
	"net"

	"github.com/jonathantorres/voy/internal/http"
)

// starts the server process and handles every request sent to it
// handles server start, restart and shutdown

const (
	name     = "localhost"
	port     = 8010
	buffSize = 1024
)

func Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", name, port))
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
	l.Close()
	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reqData := make([]byte, buffSize)
	_, err := conn.Read(reqData)
	if err != nil {
		conn.Write(http.BuildResponseBytes(http.SendServerError()))
		log.Println(err)
		return
	}
	// build the req object based on these bytes of data
	// should we return an error here?
	// or should the server just send a specific response?
	req, err := http.NewRequest(reqData)
	if err != nil {
		log.Println(err)
		// TODO: send appropriate error response
		return
	}

	// build the response string and return it
	res := http.NewResponse(req)
	_, err = conn.Write(http.BuildResponseBytes(res))
	if err != nil {
		conn.Write(http.BuildResponseBytes(http.SendServerError()))
		log.Println(err)
	}
}
