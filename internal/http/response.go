package http

import (
	"fmt"

	"github.com/jonathantorres/voy/internal/voy"
)

type Response struct {
	httpVersionMajor int
	httpVersionMinor int
	code             int
	message          string
	headers          map[string]string
	body             []byte
}

func NewResponse(req *Request) *Response {
	headers := make(map[string]string)
	body := make([]byte, 0)

	body = append(body, []byte("Hola!")...) // TODO
	addDefaultResponseHeaders(headers)
	headers["Content-Type"] = "text/html" // TODO
	res := &Response{
		httpVersionMinor: req.httpVersionMinor,
		httpVersionMajor: req.httpVersionMajor,
		code:             200,  // TODO
		message:          "OK", // TODO
		headers:          headers,
		body:             body,
	}
	return res
	// return "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nServer: voy\r\n\r\n<p>Hola!</p>"
}

func BuildResponseBytes(res *Response) []byte {
	resBytes := make([]byte, 0)
	resBytes = append(resBytes, []byte(fmt.Sprintf("HTTP/%d.%d %d %s\r\n", res.httpVersionMajor, res.httpVersionMinor, res.code, res.message))...)

	for k, v := range res.headers {
		resBytes = append(resBytes, []byte(fmt.Sprintf("%s: %s\r\n", k, v))...)
	}
	resBytes = append(resBytes, []byte("\r\n")...)
	resBytes = append(resBytes, res.body...)
	return resBytes
}

func SendServerError() *Response {
	headers := make(map[string]string)
	addDefaultResponseHeaders(headers)
	headers["Connection"] = "close"
	return &Response{
		httpVersionMajor: 1,
		httpVersionMinor: 1,
		code:             500,
		message:          "Internal Server Error",
		headers:          headers,
		body:             []byte(fmt.Sprintf("%d %s", 500, "Internal Server Error")),
	}
}

func addDefaultResponseHeaders(headers map[string]string) {
	headers["Server"] = "voy v" + voy.Version
}
