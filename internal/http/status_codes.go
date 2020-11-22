package http

import "errors"

const (
	HttpVersionMajor = 1
	HttpVersionMinor = 1
)

const (
	StatusContinue          = 100
	StatusSwichingProtocols = 101
	StatusProcessing        = 102
	StatusEarlyHints        = 103

	StatusOk                          = 200
	StatusCreated                     = 201
	StatusAccepted                    = 202
	StatusNonAuthoritativeInformation = 203
	StatusNoContent                   = 204
	StatusResetContent                = 205
	StatusPartialContent              = 206
	StatusMultiStatus                 = 207
	StatusAlreadyReported             = 208
	StatusImUsed                      = 226

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
	StatusUriTooLong                  = 414
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
	StatusHttpVersionNotSupported       = 505
	StatusVariantAlsoNegotiates         = 506
	StatusInsufficientStorage           = 507
	StatusLoopDetected                  = 508
	StatusNotExtended                   = 510
	StatusNetworkAuthenticationRequired = 511
)

var statusCodes = map[int]string{
	StatusContinue:          "Continue",
	StatusSwichingProtocols: "Swiching Protocols",
	StatusProcessing:        "Processing",
	StatusEarlyHints:        "EarlyHints",

	StatusOk:                          "OK",
	StatusCreated:                     "Created",
	StatusAccepted:                    "Accepted",
	StatusNonAuthoritativeInformation: "Non-Authoritative Information",
	StatusNoContent:                   "No Content",
	StatusResetContent:                "Reset Content",
	StatusPartialContent:              "Partial Content",
	StatusMultiStatus:                 "Multi-Status",
	StatusAlreadyReported:             "Already Reported",
	StatusImUsed:                      "IM Used",

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
	StatusUriTooLong:                  "URI Too Long ",
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
	StatusHttpVersionNotSupported:       "HTTP Version Not Supported",
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
