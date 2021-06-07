package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/textproto"
	"net/url"
	"strconv"
	"unicode"
)

var ErrInvalidRequestLine = errors.New("invalid request line")
var ErrInvalidRequestMethod = errors.New("invalid request method")
var ErrRequestBodyRequired = errors.New("request body required")

const buffSize = 1024

type Request struct {
	Method           string
	Uri              string
	HTTPVersionMajor int
	HTTPVersionMinor int
	Headers          textproto.MIMEHeader
	Body             io.Reader
	r                io.Reader
	tr               *textproto.Reader
}

func (r *Request) Parse() error {
	if r.r == nil || r.tr == nil {
		panic("a reader must be specified")
	}
	if err := r.parseRequestLine(); err != nil {
		return err
	}
	if err := r.parseRequestHeaders(); err != nil {
		return err
	}
	if r.Method == RequestMethodPost {
		// TODO: setting the body to a reader for now
		// not sure how an HTTP server should handle the body
		// of the request, maybe it sends it to whatever
		// script/service needs it!
		r.Body = r.tr.R
	}
	return nil
}

func (r *Request) parseRequestLine() error {
	b, err := r.tr.ReadLineBytes()
	if err != nil && err != io.EOF {
		return err
	}
	line := bytes.Split(b, []byte(" "))
	if len(line) != 3 {
		return ErrInvalidRequestLine
	}
	r.Method, r.Uri = string(line[0]), string(line[1])

	// validate the request method
	if err = r.validateMethod(); err != nil {
		return err
	}
	// validate the request uri
	if err = r.validateURI(); err != nil {
		return err
	}
	// TODO: validate the HTTP version,
	// it should be a valid one supported by the server
	// parse the HTTP version
	var major, minor int
	for _, char := range line[2] {
		if unicode.IsDigit(rune(char)) {
			if major == 0 {
				major, err = strconv.Atoi(string(char))
				if err != nil {
					return ErrInvalidRequestLine
				}
			} else {
				minor, err = strconv.Atoi(string(char))
				if err != nil {
					return ErrInvalidRequestLine
				}
			}
		}
	}
	r.HTTPVersionMajor, r.HTTPVersionMinor = major, minor
	// log.Printf("req line: %s\n", string(b))
	return nil
}

func (r *Request) parseRequestHeaders() error {
	h, err := r.tr.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return err
	}
	// log.Printf("headers read: %v\n", h)
	r.Headers = h
	return nil
}

func (r *Request) validateMethod() error {
	switch r.Method {
	case RequestMethodGet:
	case RequestMethodHead:
	case RequestMethodPut:
	case RequestMethodDelete:
	case RequestMethodTrace:
	case RequestMethodOptions:
	case RequestMethodConnect:
	case RequestMethodPatch:
		// everything ok, this request method is allowed
		return nil

	case RequestMethodPost:
		if r.Body == nil {
			return ErrRequestBodyRequired
		}
		return nil
	}
	return ErrInvalidRequestMethod
}

func (r *Request) validateURI() error {
	if _, err := url.ParseRequestURI(r.Uri); err != nil {
		return err
	}
	return nil
}

func NewRequest(r io.Reader) *Request {
	tr := textproto.NewReader(bufio.NewReaderSize(r, buffSize))
	return &Request{
		r:  r,
		tr: tr,
	}
}
