package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/textproto"
	"net/url"
	"regexp"
	"strconv"
	"unicode"
)

var (
	ErrInvalidRequestLine      = errors.New("invalid request line")
	ErrInvalidRequestMethod    = errors.New("invalid request method")
	ErrInvalidHTTPVersion      = errors.New("invalid http version")
	ErrRequestBodyRequired     = errors.New("request body required")
	ErrHTTPVersionNotSupported = errors.New("http version not supported")
)
var httpRegex = regexp.MustCompile(`HTTP\/\d{1}\.\d{1}`)

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
	// parse the HTTP version
	var major, minor int
	for _, char := range line[2] {
		if unicode.IsDigit(rune(char)) {
			if major == 0 {
				major, err = strconv.Atoi(string(char))
				if err != nil {
					return ErrInvalidHTTPVersion
				}
			} else {
				minor, err = strconv.Atoi(string(char))
				if err != nil {
					return ErrInvalidHTTPVersion
				}
			}
		}
	}
	r.HTTPVersionMajor, r.HTTPVersionMinor = major, minor
	// validate the HTTP version,
	// it should be a valid one supported by the server
	if err = r.validateHTTPVersion(string(line[2])); err != nil {
		return err
	}
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
	case RequestMethodGet,
		RequestMethodHead,
		RequestMethodPut,
		RequestMethodDelete,
		RequestMethodTrace,
		RequestMethodOptions,
		RequestMethodConnect,
		RequestMethodPatch:
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

func (r *Request) validateHTTPVersion(v string) error {
	if r.HTTPVersionMajor > 1 {
		return ErrHTTPVersionNotSupported
	}
	if r.HTTPVersionMajor == 0 && r.HTTPVersionMinor != 9 {
		return ErrHTTPVersionNotSupported
	}
	if r.HTTPVersionMinor > 1 {
		return ErrHTTPVersionNotSupported
	}
	if ok := httpRegex.MatchString(v); !ok {
		return ErrInvalidHTTPVersion
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
