package request

import (
	"fmt"
	"io"
	"strings"
)

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

var ERROR_MALFORMED_RLINE = fmt.Errorf("malformed request line oopsie!")
var SEPARATOR = "\r\n"

func parseRequestLine(s string) (*RequestLine, string, error) {
	idx := strings.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, s, ERROR_MALFORMED_RLINE
	}
	startLine := s[:idx]
	restOfMsg := s[idx+len(SEPARATOR):]
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, s, ERROR_MALFORMED_RLINE
	}
	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HTTPVersion:   parts[2],
	}, restOfMsg, nil
}

func (r Request) Do() (*Response, error) {
	return &Response{Body: io.NopCloser(nil)}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	requestLine, _, err := parseRequestLine(string(content))
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *requestLine,
		Headers:     nil,
		Body:        nil,
	}, nil
}
