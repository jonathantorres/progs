package server

import (
	"fmt"
	"net"
	"os"
)

// starts the server process and handles every request sent to it
// handles server start, restart and shutdown

const (
	name = "localhost"
	port = 8010
)

func Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", name, port))
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n")
			continue
		}
		go handleConn(conn)
	}
	l.Close()
	return nil
}

func handleConn(conn net.Conn) {
	msg := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nServer: voy\r\n\r\n<p>Hola!</p>"
	req := make([]byte, 512)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n")
		return
	}
	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n")
		return
	}
	conn.Close()
}
