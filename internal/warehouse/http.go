package warehouse

type HTTPResponse struct {
	HTTPVersion string // HTTP/1.1
	statuscode  int    // 200
	reason      string // ok
	body        string // content
    headers     []HTTPHeader
    endHeaders  string
}

type HTTPHeader struct {
    name    string
    value   string
}

func NewHTTPResponse() (res *HTTPResponse, err error) {
	response := &HTTPResponse{
		HTTPVersion: "HTTP/1.1",
		statuscode:  200,
		reason:      "OK \r\n",
        headers:     []HTTPHeader{},
        endHeaders:  "\r\n\r\n",
		body:        "All Sensor Data",
	}

    return response, nil
}

func (r *HTTPResponse) SetHeader(name, value string) error {
    return nil
}

func (r *HTTPResponse) SetStatusCode(code int) error {
    r.statuscode = code
    return nil
}
