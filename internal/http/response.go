package http

import (
	"fmt"
	"log"

	"github.com/jonathantorres/voy/internal/voy"
)

type Response struct {
	HTTPVersionMajor int
	HTTPVersionMinor int
	Code             int
	Message          string
	Headers          map[string]string
	Body             []byte
}

func NewResponse(code int, headers map[string]string, body []byte) *Response {
	msg, err := GetStatusCodeMessage(code)
	if err != nil {
		// TODO: handle errors better here :)
		log.Println(err)
	}
	addDefaultResponseHeaders(headers)
	res := &Response{
		HTTPVersionMinor: HTTPVersionMinor,
		HTTPVersionMajor: HTTPVersionMajor,
		Code:             code,
		Message:          msg,
		Headers:          headers,
		Body:             body,
	}
	return res
}

func BuildResponseBytes(res *Response) []byte {
	resBytes := make([]byte, 0)
	resBytes = append(resBytes, []byte(fmt.Sprintf("HTTP/%d.%d %d %s\r\n", res.HTTPVersionMajor, res.HTTPVersionMinor, res.Code, res.Message))...)

	for k, v := range res.Headers {
		resBytes = append(resBytes, []byte(fmt.Sprintf("%s: %s\r\n", k, v))...)
	}
	resBytes = append(resBytes, []byte("\r\n")...)
	resBytes = append(resBytes, res.Body...)
	return resBytes
}

func SendErrorResponse(code int, msg string) *Response {
	headers := make(map[string]string)
	addDefaultResponseHeaders(headers)
	headers["Content-Type"] = "text/html"
	headers["Connection"] = "close"
	return &Response{
		HTTPVersionMajor: HTTPVersionMinor,
		HTTPVersionMinor: HTTPVersionMajor,
		Code:             code,
		Message:          msg,
		Headers:          headers,
		Body:             []byte(fmt.Sprintf("%d %s", code, msg)),
	}
}

func addDefaultResponseHeaders(headers map[string]string) {
	headers["Server"] = "voy v" + voy.Version
}
