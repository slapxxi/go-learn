package request

import "io"

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HTTPVersion   string
}

type Response struct {
	Body io.ReadCloser
}

func (r Request) Do() (*Response, error) {
	return &Response{Body: io.NopCloser(nil)}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	return nil, nil
}
