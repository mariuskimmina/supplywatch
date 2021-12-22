package http

import (
	"fmt"
	"strconv"
	"strings"
)

type Response struct {
	HTTPVersion string // HTTP/1.1
	statuscode  int    // 200
	reason      string // ok
	headers     string // key:value\r\n
	body        []byte // content
}

type ResponseOption func(*Response) error

func WithHTTPVersion(version string) func(*Response) error {
    return func(res *Response) error {
        res.HTTPVersion = version
        return nil
    }
}

func WithStatusCode(code int) func(*Response) error {
    return func(res *Response) error {
        res.statuscode = code
        return nil
    }
}

func WithReason(reason string) func(*Response) error {
    return func(res *Response) error {
        res.reason = reason
        return nil
    }
}

func WithHeaders(headers string) func(*Response) error {
    return func(res *Response) error {
        res.headers = headers
        return nil
    }
}

func WithBody(body []byte) func(*Response) error {
    return func(res *Response) error {
        res.body = body
        return nil
    }
}

const (
    HTTPResponseDefaultStatusCode int = 200
    HTTPResponseDefaultHTTPVersion string = "HTTP/1.1"
    HTTPResponseDefaultReason string = "OK\r\n"
    HTTPResponseDefaultHeaders string = "\r\n"
)

var (
    HTTPResponseDefaultBody []byte = []byte{}
)


func NewResponse(opts ...ResponseOption) (res *Response, err error) {
	response := &Response{
		HTTPVersion: HTTPResponseDefaultHTTPVersion,
		statuscode:  HTTPResponseDefaultStatusCode,
		reason:      HTTPResponseDefaultReason,
		headers:     HTTPResponseDefaultHeaders,
		body:        HTTPResponseDefaultBody,
	}

    for _, opt := range opts {
        if err := opt(response); err != nil {
            return &Response{}, fmt.Errorf("option failed %w", err)
        }
    }

	return response, nil
}

func ResponseToBytes(res *Response) (response []byte, err error) {
	var parts []string
	for _, s := range []string{res.HTTPVersion, strconv.Itoa(res.statuscode), res.reason, res.headers, string(res.body)} {
		parts = append(parts, s)
	}
	stringResponse := strings.Join(parts, " ")
	response = []byte(stringResponse)
	return response, nil
}

// SetHeader first removes the newline at the end of headers
// it then creates a new header, inserts the header into headers
// and appends a newline again to end the headers section
func (r *Response) SetHeader(name, value string) error {
	r.headers = r.headers[:len(r.headers)-2]
	header := NewHeader(name, value)
	r.headers += header
	endOfHeaders := "\r\n"
	r.headers += endOfHeaders
	return nil
}

func (r *Response) SetBody(body []byte) error {
	r.body = body
	return nil
}

func (r *Response) SetStatusCode(code int) error {
	r.statuscode = code
	return nil
}

func (r *Response) SetReason(reason string) error {
	reason = reason + "\r\n"
	r.reason = reason
	return nil
}
