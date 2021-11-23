package warehouse

import (
	"strconv"
	"strings"
)

type HTTPResponse struct {
	HTTPVersion string // HTTP/1.1
	statuscode  int    // 200
	reason      string // ok
	headers     string // key:value\r\n
	body        []byte // content
}

type HTTPRequest struct {
	method  string
	path    string
	version string
	query   string
}

type HTTPHeader struct {
	string
}

func ResponseToBytes(res *HTTPResponse) (response []byte, err error) {
	var parts []string
	for _, s := range []string{res.HTTPVersion, strconv.Itoa(res.statuscode), res.reason, res.headers, string(res.body)} {
		parts = append(parts, s)
	}
	stringResponse := strings.Join(parts, "")
	response = []byte(stringResponse)
	return response, nil
}

func NewHTTPResponse() (res *HTTPResponse, err error) {
	response := &HTTPResponse{
		HTTPVersion: "HTTP/1.1",
		statuscode:  200,
		reason:      "OK\r\n",
		headers:     "\r\n",
		body:        []byte{},
	}

	return response, nil
}

func NewHTTPHeader(name, value string) string {
	endOfHeader := "\r\n"
	header := name + ":" + value + endOfHeader
	return header
}

// SetHeader first removes the newline at the end of headers
// it then creates a new header, inserts the header into headers
// and appends a newline again to end the headers section
func (r *HTTPResponse) SetHeader(name, value string) error {
	r.headers = r.headers[:len(r.headers)-2]
	header := NewHTTPHeader(name, value)
	r.headers += header
	endOfHeaders := "\r\n"
	r.headers += endOfHeaders
	return nil
}

func (r *HTTPResponse) SetBody(body []byte) error {
	r.body = body
	return nil
}

func (r *HTTPResponse) SetStatusCode(code int) error {
	r.statuscode = code
	return nil
}
