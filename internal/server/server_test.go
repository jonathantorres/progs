package server

import (
	"testing"
)

var req string = "GET / HTTP/1.1\r\nHost: www.example.com\r\n\r\n"

func TestRequestLine(t *testing.T) {
	req := newRequest([]byte(req))
	if req.method != "GET" {
		t.Errorf("req method should be GET")
	}
	if req.uri != "/" {
		t.Errorf("req uri should be /")
	}
	if req.httpVersionMajor != 1 {
		t.Errorf("req major version should be 1")
	}
	if req.httpVersionMinor != 1 {
		t.Errorf("req minor version should be 1")
	}
}
