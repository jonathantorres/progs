package http

import (
	"fmt"
	"log"

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

func NewResponse(code int, headers map[string]string, body []byte) *Response {
	msg, err := GetStatusCodeMessage(code)
	if err != nil {
		// TODO: handle errors better here :)
		log.Println(err)
	}
	addDefaultResponseHeaders(headers)
	res := &Response{
		httpVersionMinor: HTTPVersionMinor,
		httpVersionMajor: HTTPVersionMajor,
		code:             code,
		message:          msg,
		headers:          headers,
		body:             body,
	}
	return res
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

func SendErrorResponse(code int, msg string) *Response {
	headers := make(map[string]string)
	addDefaultResponseHeaders(headers)
	headers["Content-Type"] = "text/html"
	headers["Connection"] = "close"
	return &Response{
		httpVersionMajor: HTTPVersionMinor,
		httpVersionMinor: HTTPVersionMajor,
		code:             code,
		message:          msg,
		headers:          headers,
		body:             []byte(fmt.Sprintf("%d %s", code, msg)),
	}
}

func addDefaultResponseHeaders(headers map[string]string) {
	headers["Server"] = "voy v" + voy.Version
}
