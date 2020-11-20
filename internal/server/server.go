package server

import (
	"errors"
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

var (
	serverErrResponseCode = 500
	clientErrResponseCode = 400
	serverErrResponseMsg  = "Internal Server Error"
	clientErrResponseMsg  = "Bad Request"
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
		conn.Write(http.BuildResponseBytes(http.SendErrorResponse(serverErrResponseCode, serverErrResponseMsg)))
		log.Println(err)
		return
	}
	req, err := http.NewRequest(reqData)
	if err != nil {
		if errors.Is(err, http.ErrInvalidRequestLine) {
			conn.Write(http.BuildResponseBytes(http.SendErrorResponse(clientErrResponseCode, clientErrResponseMsg)))
		} else {
			conn.Write(http.BuildResponseBytes(http.SendErrorResponse(serverErrResponseCode, serverErrResponseMsg)))
		}
		log.Println(err)
		return
	}

	res := http.NewResponse(req)
	_, err = conn.Write(http.BuildResponseBytes(res))
	if err != nil {
		conn.Write(http.BuildResponseBytes(http.SendErrorResponse(serverErrResponseCode, serverErrResponseMsg)))
		log.Println(err)
	}
}
