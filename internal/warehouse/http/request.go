package http

import "fmt"


type Request struct {
	Method  string
	Path    string
	Version string
	Query   string
}

type RequestOption func(*Request) error

func WithMethod(method string) func(*Request) error {
    return func(req *Request) error {
        req.Method = method
        return nil
    }
}

func WithPath(path string) func(*Request) error {
    return func(req *Request) error {
        req.Path = path
        return nil
    }
}

func WithVersion(version string) func(*Request) error {
    return func(req *Request) error {
        req.Version = version
        return nil
    }
}

func WithQuery(query string) func(*Request) error {
    return func(req *Request) error {
        req.Query = query
        return nil
    }
}

const (
    HTTPRequestDefaultMethod string = ""
    HTTPRequestDefaultPath string = ""
    HTTPRequestDefaultVersion string = ""
    HTTPRequestDefaultQuery string = ""
)

var (
    HTTPRequestDefaultBody []byte = []byte{}
)


func NewRequest(opts ...RequestOption) (res *Request, err error) {
	request := &Request{
		Method: HTTPRequestDefaultMethod,
		Path:  HTTPRequestDefaultPath,
		Version:      HTTPRequestDefaultVersion,
		Query:     HTTPRequestDefaultQuery,
	}

    for _, opt := range opts {
        if err := opt(request); err != nil {
            return &Request{}, fmt.Errorf("option failed %w", err)
        }
    }

	return request, nil
}
