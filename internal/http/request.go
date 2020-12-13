package http

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"unicode"
)

type Request struct {
	Method           string
	Uri              string
	httpVersionMajor int
	httpVersionMinor int
	Headers          map[string]string
	Body             []byte
}

var ErrInvalidRequestLine = errors.New("invalid request line")

func NewRequest(reqData []byte) (*Request, error) {
	method, uri, major, minor, err := parseRequestLine(reqData)
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
		httpVersionMajor: major,
		httpVersionMinor: minor,
		Headers:          headers,
		Body:             body,
	}
	return req, nil
}

func parseRequestLine(reqData []byte) (string, string, int, int, error) {
	var method, uri string
	var major, minor int
	var tok []byte
	r := bytes.NewReader(reqData)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tok = scanner.Bytes()
		break // only read the first line
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
		parts := bytes.Split(tok, []byte(":"))
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
