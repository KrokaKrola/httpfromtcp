package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type Request struct {
	RequestLine RequestLine
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(rawBytes)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request line")
	}

	requestLine, err := requestLineFromString(string(data[:idx]))

	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
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
