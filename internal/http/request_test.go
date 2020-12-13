package http

import (
	"reflect"
	"testing"
)

var pay1 = `GET / HTTP/1.1
Host: www.example.com
`
var pay1Headers = map[string]string{
	"Host": "www.example.com",
}
var pay1Body []byte = nil

var pay2 = `GET /foo/bar HTTP/1.1
Host: www.example.com
Server: voy v0.1.0
Connection: close
`
var pay2Headers = map[string]string{
	"Host":       "www.example.com",
	"Server":     "voy v0.1.0",
	"Connection": "close",
}
var pay2Body []byte = nil

var pay3 = `POST /user/create HTTP/1.1
Host: www.example.com
Server: voy v0.1.0
Connection: close
Content-Length: 41

user=foo&password=bar&email=test@test.com
`
var pay3Headers = map[string]string{
	"Host":           "www.example.com",
	"Server":         "voy v0.1.0",
	"Connection":     "close",
	"Content-Length": "41",
}
var pay3Body = []byte("user=foo&password=bar&email=test@test.com")

var pay4 = `GET / HTTP/1.1
Host: localhost:8021
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
DNT: 1
Connection: keep-alive
Upgrade-Insecure-Requests: 1
Pragma: no-cache
Cache-Control: no-cache
`
var pay4Headers = map[string]string{
	"Host":                      "localhost:8021",
	"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Language":           "en-US,en;q=0.5",
	"Accept-Encoding":           "gzip, deflate",
	"DNT":                       "1",
	"Connection":                "keep-alive",
	"Upgrade-Insecure-Requests": "1",
	"Pragma":                    "no-cache",
	"Cache-Control":             "no-cache",
}
var pay4Body []byte = nil

var cases = []struct {
	payload          string
	method           string
	uri              string
	httpVersionMajor int
	httpVersionMinor int
	headers          map[string]string
	body             []byte
}{
	{pay1, "GET", "/", 1, 1, pay1Headers, pay1Body},
	{pay2, "GET", "/foo/bar", 1, 1, pay2Headers, pay2Body},
	{pay3, "POST", "/user/create", 1, 1, pay3Headers, pay3Body},
	{pay4, "GET", "/", 1, 1, pay4Headers, pay4Body},
}

func TestRequestLine(t *testing.T) {
	for _, c := range cases {
		req, err := NewRequest([]byte(c.payload))
		if err != nil {
			t.Errorf(err.Error())
		}
		if req.Method != c.method {
			t.Errorf("req method is %s but it should be %s", req.Method, c.method)
		}
		if req.Uri != c.uri {
			t.Errorf("req uri is %s and it should be %s", req.Uri, c.uri)
		}
		if req.httpVersionMajor != c.httpVersionMajor {
			t.Errorf("req major version is %d but it should be %d", req.httpVersionMajor, c.httpVersionMajor)
		}
		if req.httpVersionMinor != c.httpVersionMinor {
			t.Errorf("req minor version is %d but it should be %d", req.httpVersionMinor, c.httpVersionMinor)
		}
	}
}

func TestParsingOfHeaders(t *testing.T) {
	for i, c := range cases {
		req, err := NewRequest([]byte(c.payload))
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reflect.DeepEqual(req.Headers, c.headers) {
			t.Errorf("headers from payload#%d are not equal", i+1)
		}
	}
}

func TestParsingOfBody(t *testing.T) {
	for i, c := range cases {
		req, err := NewRequest([]byte(c.payload))
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reflect.DeepEqual(req.Body, c.body) {
			t.Errorf("request body from payload#%d is not equal", i+1)
		}
	}
}
