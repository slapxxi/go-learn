package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

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

func (r *RequestLine) ValidHTTP() bool {
	return r.HTTPVersion == "1.1"
}

type Response struct {
	Body io.ReadCloser
}

var ERROR_MALFORMED_RLINE = fmt.Errorf("malformed request line oopsie!")
var ERROR_UNSUPPORTED_VERSION = fmt.Errorf("unsupported http version")
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
	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" {
		return nil, s, ERROR_MALFORMED_RLINE
	}
	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HTTPVersion:   httpParts[1],
	}
	if !rl.ValidHTTP() {
		return nil, s, ERROR_UNSUPPORTED_VERSION
	}
	return rl, restOfMsg, nil
}

func (r Request) Do() (*Response, error) {
	return &Response{Body: io.NopCloser(nil)}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Unable to read contents"), err)
	}
	str := string(data)
	rl, str, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *rl,
		Headers:     nil,
		Body:        nil,
	}, nil
}
