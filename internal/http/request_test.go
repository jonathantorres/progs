package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"
	"reflect"
	"testing"
)

var cases = []struct {
	payload          string
	method           string
	uri              string
	httpVersionMajor int
	httpVersionMinor int
	headers          map[string][]string
}{
	{
		"get_payload1", "GET", "/", 1, 1,
		map[string][]string{
			"Host": {"www.example.com"},
		},
	},
	{
		"get_payload2", "GET", "/foo/bar", 1, 1,
		map[string][]string{
			"Host":       {"www.example.com"},
			"Server":     {"voy v0.1.0"},
			"Connection": {"close"},
		},
	},
	{
		"post_payload3", "POST", "/user/create", 1, 1,
		map[string][]string{
			"Host":           {"www.example.com"},
			"Server":         {"voy v0.1.0"},
			"Connection":     {"close"},
			"Content-Length": {"41"},
		},
	},
	{
		"get_payload4", "GET", "/", 1, 1,
		map[string][]string{
			"Host":                      {"localhost:8021"},
			"User-Agent":                {"Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0"},
			"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			"Accept-Language":           {"en-US,en;q=0.5"},
			"Accept-Encoding":           {"gzip, deflate"},
			"Dnt":                       {"1"},
			"Connection":                {"keep-alive"},
			"Upgrade-Insecure-Requests": {"1"},
			"Pragma":                    {"no-cache"},
			"Cache-Control":             {"no-cache"},
		},
	},
}

func TestRequestLine(t *testing.T) {
	for _, c := range cases {
		r, err := loadRequestPayload(c.payload)
		if err != nil {
			t.Fatalf(err.Error())
		}
		req := NewRequest(r)
		err = req.Parse()
		if err != nil && err != io.EOF {
			t.Fatalf(err.Error())
		}
		if req.Method != c.method {
			t.Fatalf("req method is %s but it should be %s", req.Method, c.method)
		}
		if req.Uri != c.uri {
			t.Fatalf("req uri is %s and it should be %s", req.Uri, c.uri)
		}
		if req.HTTPVersionMajor != c.httpVersionMajor {
			t.Fatalf("req major version is %d but it should be %d", req.HTTPVersionMajor, c.httpVersionMajor)
		}
		if req.HTTPVersionMinor != c.httpVersionMinor {
			t.Fatalf("req minor version is %d but it should be %d", req.HTTPVersionMinor, c.httpVersionMinor)
		}
	}
}

func TestParsingOfHeaders(t *testing.T) {
	for i, c := range cases {
		r, err := loadRequestPayload(c.payload)
		if err != nil {
			t.Fatalf(err.Error())
		}
		req := NewRequest(r)
		err = req.Parse()
		if err != nil && err != io.EOF {
			t.Fatalf(err.Error())
		}
		if !reflect.DeepEqual(req.Headers, textproto.MIMEHeader(c.headers)) {
			t.Fatalf("headers from payload#%d are not equal, got %v but want %v", i+1, req.Headers, c.headers)
		}
	}
}

func TestParsingOfBody(t *testing.T) {
	// we probably don't need this anymore
	t.SkipNow()
}

func loadRequestPayload(payload string) (io.Reader, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.txt", payload))
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)
	return r, nil
}
