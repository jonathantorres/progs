package main

import (
	"errors"
	"fmt"
	"log"
)

const (
	HTTPVersionMajor = 1
	HTTPVersionMinor = 1
)

const (
	StatusContinue           = 100
	StatusSwitchingProtocols = 101
	StatusProcessing         = 102
	StatusEarlyHints         = 103

	StatusOk                          = 200
	StatusCreated                     = 201
	StatusAccepted                    = 202
	StatusNonAuthoritativeInformation = 203
	StatusNoContent                   = 204
	StatusResetContent                = 205
	StatusPartialContent              = 206
	StatusMultiStatus                 = 207
	StatusAlreadyReported             = 208
	StatusIMUsed                      = 226

	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	StatusSwitchProxy       = 306
	StatusTemporaryRedirect = 307
	StatusPermanentRedirect = 308

	StatusBadRequest                  = 400
	StatusUnauthorized                = 401
	StatusPaymentRequired             = 402
	StatusForbidden                   = 403
	StatusNotFound                    = 404
	StatusMethodNotAllowed            = 405
	StatusNotAcceptable               = 406
	StatusProxyAuthenticationRequired = 407
	StatusRequestTimeout              = 408
	StatusConflict                    = 409
	StatusGone                        = 410
	StatusLengthRequired              = 411
	StatusPreconditionFailed          = 412
	StatusPayloadTooLarge             = 413
	StatusURITooLong                  = 414
	StatusUnsupportedMediaType        = 415
	StatusRangeNotSatisfiable         = 416
	StatusExpectationFailed           = 417
	StatusImATeapot                   = 418
	StatusMisdirectedRequest          = 421
	StatusUnprocessableEntity         = 422
	StatusLocked                      = 423
	StatusFailedDependency            = 424
	StatusTooEarly                    = 425
	StatusUpgradeRequired             = 426
	StatusPreconditionRequired        = 428
	StatusTooManyRequests             = 429
	StatusRequestHeaderFieldsTooLarge = 431
	StatusUnavailableForLegalReasons  = 451

	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusVariantAlsoNegotiates         = 506
	StatusInsufficientStorage           = 507
	StatusLoopDetected                  = 508
	StatusNotExtended                   = 510
	StatusNetworkAuthenticationRequired = 511
)

var statusCodes = map[int]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",
	StatusProcessing:         "Processing",
	StatusEarlyHints:         "EarlyHints",

	StatusOk:                          "OK",
	StatusCreated:                     "Created",
	StatusAccepted:                    "Accepted",
	StatusNonAuthoritativeInformation: "Non-Authoritative Information",
	StatusNoContent:                   "No Content",
	StatusResetContent:                "Reset Content",
	StatusPartialContent:              "Partial Content",
	StatusMultiStatus:                 "Multi-Status",
	StatusAlreadyReported:             "Already Reported",
	StatusIMUsed:                      "IM Used",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusSwitchProxy:       "Switch Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",
	StatusPermanentRedirect: "Permanent Redirect",

	StatusBadRequest:                  "Bad Request",
	StatusUnauthorized:                "Unauthorized",
	StatusPaymentRequired:             "Payment Required",
	StatusForbidden:                   "Forbidden",
	StatusNotFound:                    "Not Found",
	StatusMethodNotAllowed:            "Method Not Allowed",
	StatusNotAcceptable:               "Not Acceptable",
	StatusProxyAuthenticationRequired: "Proxy Authentication Required",
	StatusRequestTimeout:              "Request Timeout",
	StatusConflict:                    "Conflict",
	StatusGone:                        "Gone",
	StatusLengthRequired:              "Length Required",
	StatusPreconditionFailed:          "Precondition Failed",
	StatusPayloadTooLarge:             "Payload Too Large",
	StatusURITooLong:                  "URI Too Long ",
	StatusUnsupportedMediaType:        "Unsupported Media Type",
	StatusRangeNotSatisfiable:         "Range Not Satisfiable ",
	StatusExpectationFailed:           "Expectation Failed",
	StatusImATeapot:                   "I'm a teapot",
	StatusMisdirectedRequest:          "Misdirected Request",
	StatusUnprocessableEntity:         "Unprocessable Entity",
	StatusLocked:                      "Locked",
	StatusFailedDependency:            "Failed Dependency",
	StatusTooEarly:                    "Too Early",
	StatusUpgradeRequired:             "Upgrade Required",
	StatusPreconditionRequired:        "Precondition Required",
	StatusTooManyRequests:             "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge: "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:  "Unavailable For Legal Reasons",

	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	StatusInsufficientStorage:           "Insufficient Storage",
	StatusLoopDetected:                  "Loop Detected",
	StatusNotExtended:                   "Not Extended ",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

func GetStatusCodeMessage(statusCode int) (string, error) {
	for code, statusMsg := range statusCodes {
		if code == statusCode {
			return statusMsg, nil
		}
	}
	return "", errors.New("status code not found")
}

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
	headers["Server"] = "httpd v" + Version
}
