package server

import (
	"testing"
)

var pay1 = `GET / HTTP/1.1
Host: www.example.com
`
var pay2 = `GET /foo/bar HTTP/1.1
Host: www.example.com
Server: voy v0.1.0
Connection: close
`

var cases = []struct {
	payload          string
	method           string
	uri              string
	httpVersionMajor int
	httpVersionMinor int
}{
	{pay1, "GET", "/", 1, 1},
	{pay2, "GET", "/foo/bar", 1, 1},
}

func TestRequestLine(t *testing.T) {
	for _, c := range cases {
		req := newRequest([]byte(c.payload))
		if req.method != c.method {
			t.Errorf("req method is %s but it should be %s", req.method, c.method)
		}
		if req.uri != c.uri {
			t.Errorf("req uri is %s and it should be %s", req.uri, c.uri)
		}
		if req.httpVersionMajor != c.httpVersionMajor {
			t.Errorf("req major version is %d but it should be %d", req.httpVersionMajor, c.httpVersionMajor)
		}
		if req.httpVersionMinor != c.httpVersionMinor {
			t.Errorf("req minor version is %d but it should be %d", req.httpVersionMinor, c.httpVersionMinor)
		}
	}
}
