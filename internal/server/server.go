package server

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"unicode"
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

func parseHeaders(reqData []byte) (map[string]string, error) {
	return nil, nil
}

func newRequest(reqData []byte) *Request {
	method, uri, major, minor, err := parseRequestLine(reqData)
	if err != nil {
		return nil
	}
	headers, err := parseHeaders(reqData)
	if err != nil {
		return nil
	}
	// TODO: parse body (if any)
	return &Request{
		method:           method,
		uri:              uri,
		httpVersionMajor: major,
		httpVersionMinor: minor,
		headers:          headers,
	}
}

func parseRequestLine(reqData []byte) (string, string, int, int, error) {
	var method, uri string
	var major, minor int
	var tok []byte
	r := bytes.NewReader(reqData)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tok = scanner.Bytes()
		break // only read the first line
	}
	if scanner.Err() != nil {
		return "", "", 0, 0, scanner.Err()
	}

	parts := bytes.Split(tok, []byte{byte(' ')})
	if len(parts) != 3 {
		return "", "", 0, 0, errors.New("invalid request line")
	}
	method = string(parts[0])
	uri = string(parts[1])
	for _, char := range parts[2] {
		if unicode.IsDigit(rune(char)) {
			if major == 0 {
				major, _ = strconv.Atoi(string(char))
			} else {
				minor, _ = strconv.Atoi(string(char))
			}
		}
	}
	return method, uri, major, minor, nil
}

func newResponse(req *Request) string {
	return "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nServer: voy\r\n\r\n<p>Hola!</p>"
}

type Request struct {
	method           string
	uri              string
	httpVersionMajor int
	httpVersionMinor int
	headers          map[string]string
	body             []byte
}

type Response struct {
	httpVersionMajor int
	httpVersionMinor int
	code             int
	message          string
	headers          map[string]string
	body             []byte
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
