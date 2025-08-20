package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	GET     = "GET"
	POST    = "POST"
	PATCH   = "PATCH"
	PUT     = "PUT"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
)

var validRequestMethodsMap = map[string]bool{
	GET:     true,
	POST:    true,
	PATCH:   true,
	PUT:     true,
	DELETE:  true,
	OPTIONS: true,
}

func isAllCapitalLetters(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}

	return true
}

func parseRequestLine(s string) (*RequestLine, error) {
	fmt.Println("string", s)
	requestLineParts := strings.Split(s, " ")

	fmt.Println("requestLineParts", requestLineParts)

	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", s)
	}

	if !isAllCapitalLetters(requestLineParts[0]) {
		return nil, fmt.Errorf("invalid method name: %s", requestLineParts[0])
	}

	if !validRequestMethodsMap[requestLineParts[0]] {
		return nil, fmt.Errorf("not supported method: %s", requestLineParts[0])
	}

	httpVersion := strings.Split(requestLineParts[2], "/")

	if len(httpVersion) != 2 {
		return nil, fmt.Errorf("invalid http version: %s", httpVersion)
	}

	if httpVersion[0] != "HTTP" {
		return nil, fmt.Errorf("invalid protocol: %s", httpVersion[0])
	}

	return &RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersion[1],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	res, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestParts := strings.Split(string(res), "\r\n")

	fmt.Println("request parts", requestParts)

	if len(requestParts) == 0 {
		return nil, errors.New("invalid request")
	}

	requestLine, err := parseRequestLine(requestParts[0])

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}
