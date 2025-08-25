package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const crlf = ("\r\n")

const separator = ":"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfPos := bytes.Index(data, []byte(crlf))

	if crlfPos == -1 {
		return 0, false, nil
	}

	if crlfPos == 0 {
		return 2, true, nil
	}

	rawHeader := data[:crlfPos]

	headerPart, valuePart, foundSeparator := bytes.Cut(rawHeader, []byte(separator))

	if !foundSeparator {
		return 0, false, fmt.Errorf("invalid header")
	}

	headerName := bytes.TrimLeft(headerPart, " ")

	if len(headerName) == 0 {
		return 0, false, fmt.Errorf("invalid header name")
	}

	if bytes.HasSuffix(headerName, []byte(" ")) {
		return 0, false, fmt.Errorf("invalid header name")
	}

	if !h.isValidByteSlice(headerName) {
		return 0, false, fmt.Errorf("header contains invalid characters")
	}

	headerValueSlice := bytes.TrimSpace(valuePart)

	if len(headerValueSlice) == 0 {
		return 0, false, fmt.Errorf("invalid header value")
	}

	headerKey := string(headerName)
	headerValue := string(headerValueSlice)

	storedValue, exists := h.Get(headerKey)

	if exists {
		headerValue = fmt.Sprintf("%s, %s", storedValue, headerValue)
	}

	h.Set(headerKey, headerValue)

	return crlfPos + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, exists := h[key]

	return value, exists
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

var validBytePattern = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_\` + "`" + `|~]*$`)

// IsValidByteSlice checks if the given byte slice contains only allowed characters:
// - Uppercase letters: A-Z
// - Lowercase letters: a-z
// - Digits: 0-9
// - Special characters: !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~
func (h Headers) isValidByteSlice(data []byte) bool {
	return validBytePattern.Match(data)
}
