package http

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"unicode"
)

type Request struct {
	Method             string
	Uri                string
	HTTPVersionMajor   int
	HTTPVersionMinor   int
	Headers            map[string]string
	Body               []byte
	DoneReading        bool
	LineIsRead         bool
	HeadersAreRead     bool
	BodyIsRead         bool
	totalBodyBytesRead int
}

var ErrInvalidRequestLine = errors.New("invalid request line")

// TODO: In here we are always assuming that the buffer
// contains the entire request line
// it's possible that the request uri is large enough
// that it may not fit here, add changes to account for that
func (r *Request) ReadLine(reqData *[]byte) error {
	method, uri, major, minor, err := parseRequestLine(reqData)
	if err != nil {
		return err
	}
	r.Method = method
	r.Uri = uri
	r.HTTPVersionMajor = major
	r.HTTPVersionMinor = minor
	r.LineIsRead = true

	if r.Method == RequestMethodGet || r.Method == RequestMethodHead {
		r.BodyIsRead = true
	}
	return nil
}

// TODO: remove this function
func parseHeaders(reqData []byte) (map[string]string, error) {
	headers := make(map[string]string)
	var tok []byte
	r := bytes.NewReader(reqData)
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		tok = scanner.Bytes()
		if i == 0 {
			i++
			continue // skip the request line
		}
		if len(tok) == 0 {
			// we found an empty line, the headers end here
			break
		}
		parts := bytes.SplitN(tok, []byte(":"), 2)
		if len(parts) == 2 {
			key := string(bytes.TrimSpace(parts[0]))
			val := string(bytes.TrimSpace(parts[1]))
			headers[key] = val
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		return nil, scanner.Err()
	}
	return headers, nil
}

func (r *Request) ReadHeaders(reqData *[]byte) error {
	if !r.LineIsRead {
		return nil
	}
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	var tok []byte
	var headersLen int
	scanner := bufio.NewScanner(bytes.NewReader(*reqData))
	for scanner.Scan() {
		tok = scanner.Bytes()
		if len(tok) == 0 {
			// we found an empty line, the headers end here
			r.HeadersAreRead = true
			*reqData = (*reqData)[(headersLen + 2):] // discard the headers (if any) + the \r\n at the end of the headers
			break
		}
		headersLen += len(tok) + 2 // add the \r\n for each header line
		parts := bytes.SplitN(tok, []byte(":"), 2)
		if len(parts) == 2 {
			key := string(bytes.TrimSpace(parts[0]))
			val := string(bytes.TrimSpace(parts[1]))
			r.Headers[key] = val
		}
	}
	if err := scanner.Err(); err != nil {
		return scanner.Err()
	}
	if r.Method == RequestMethodGet || r.Method == RequestMethodHead {
		r.HeadersAreRead = true
		r.DoneReading = true
	}
	return nil
}

// TODO: remove this function
func parseBody(reqData []byte) ([]byte, error) {
	var body []byte = nil
	var tok []byte
	r := bytes.NewReader(reqData)
	scanner := bufio.NewScanner(r)
	i := 0
	foundBody := false
	for scanner.Scan() {
		tok = scanner.Bytes()
		if i == 0 {
			i++
			continue // skip the request line
		}
		if len(tok) == 0 {
			// we found an empty line, the headers end here
			// and the body starts
			foundBody = true
			continue
		}
		if foundBody {
			// this line is part of the body
			if body == nil {
				body = make([]byte, 0)
			}
			body = append(body, tok...)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return body, nil
}

func (r *Request) ReadBody(reqData *[]byte, bytesRead int) error {
	if !r.LineIsRead || !r.HeadersAreRead || r.BodyIsRead || r.DoneReading {
		return nil
	}
	var contentLen int
	if contentLenVal, ok := r.Headers["Content-Length"]; ok {
		contentLen, _ = strconv.Atoi(contentLenVal)
	}
	if r.Body == nil {
		r.Body = make([]byte, 0, contentLen)
	}
	r.Body = append(r.Body, (*reqData)...)
	r.totalBodyBytesRead += bytesRead

	if r.totalBodyBytesRead >= contentLen {
		r.BodyIsRead = true
		r.DoneReading = true
	}
	return nil
}

func NewRequest(reqData []byte) (*Request, error) {
	method, uri, major, minor, err := parseRequestLine(&reqData)
	if err != nil {
		return nil, err
	}
	headers, err := parseHeaders(reqData)
	if err != nil {
		return nil, err
	}
	body, err := parseBody(reqData)
	if err != nil {
		return nil, err
	}

	req := &Request{
		Method:           method,
		Uri:              uri,
		HTTPVersionMajor: major,
		HTTPVersionMinor: minor,
		Headers:          headers,
		Body:             body,
	}
	return req, nil
}

func parseRequestLine(reqData *[]byte) (string, string, int, int, error) {
	var method, uri string
	var major, minor int
	var tok []byte
	r := bytes.NewReader(*reqData)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tok = scanner.Bytes()
		*reqData = (*reqData)[len(tok)+2:] // discard the request line part + \r\n
		break                              // only read the first line
	}
	if err := scanner.Err(); err != nil {
		return "", "", 0, 0, err
	}
	parts := bytes.Split(tok, []byte{byte(' ')})
	if len(parts) != 3 {
		return "", "", 0, 0, ErrInvalidRequestLine
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
