package server

import (
	"fmt"
	"log"
	"net"
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

func newRequest(reqData []byte) string {
	return "request"
}

func newResponse(req string) string {
	return "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nServer: voy\r\n\r\n<p>Hola!</p>"
}

func handleConn(conn net.Conn) {
	reqData := make([]byte, buffSize)
	_, err := conn.Read(reqData)
	if err != nil {
		log.Fatal(err)
	}
	// build the req object based on these bytes of data
	// should we return an error here?
	// or should the server just send a specific response?
	req := newRequest(reqData)

	// build the response string and return it
	res := newResponse(req)
	_, err = conn.Write([]byte(res))
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
