package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
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
	state       parserState
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateDone
		case StateError:
			return 0, ERROR_RERRORSTATE
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) isDone() bool {
	return r.state == StateDone
}

func (r *Request) isError() bool {
	return r.state == StateError
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
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
var ERROR_RERRORSTATE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(s []byte) (*RequestLine, int, error) {
	idx := bytes.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, 0, ERROR_MALFORMED_RLINE
	}
	startLine := s[:idx]
	read := idx + len(SEPARATOR)
	parts := bytes.Split(startLine, []byte(""))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_RLINE
	}
	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" {
		return nil, 0, ERROR_MALFORMED_RLINE
	}
	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HTTPVersion:   string(httpParts[1]),
	}
	if !rl.ValidHTTP() {
		return nil, 0, ERROR_UNSUPPORTED_VERSION
	}
	return rl, read, nil
}

func (r Request) Do() (*Response, error) {
	return &Response{Body: io.NopCloser(nil)}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.isDone() || !request.isError() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}
