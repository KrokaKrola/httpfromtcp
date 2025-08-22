package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type requestState int

const (
	requestStateDone requestState = iota
	requestStateInitialized
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type Request struct {
	RequestLine RequestLine
	state       requestState
}

const crlf = "\r\n"
const bufSize = 1024

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufSize)
	readCounter := 0

	request := &Request{
		state: requestStateInitialized,
	}

	for request.state != requestStateDone {
		if readCounter >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readCounter:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = requestStateDone
				break
			}
			return nil, errors.Join(fmt.Errorf("error reading stream of data"), err)
		}

		readCounter += numBytesRead

		parsedCounter, err := request.parse(buf[:readCounter])

		if err != nil {
			return nil, errors.Join(fmt.Errorf("error parsing stream of data"), err)
		}

		copy(buf, buf[parsedCounter:])
		readCounter -= parsedCounter
	}

	return request, nil
}

func (r *Request) parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLine, err := r.requestLineFromString(string(data[:idx]))

	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + 2, nil
}

func (r *Request) requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", str)
	}

	method := parts[0]

	if !isAllCapitalLetters(method) {
		return nil, fmt.Errorf("invalid method name: %s", method)
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, num, err := r.parseRequestLine(data)
		// actual error happend
		if err != nil {
			return 0, err
		}

		// parsing is not done
		if num == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = requestStateDone

		return num, nil
	case requestStateDone:
		return 0, fmt.Errorf("invalid state of parser")
	default:
		return 0, fmt.Errorf("invalid state of parser")
	}
}
