package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const crlf = ("\r\n")

const separator = ":"

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
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

	h.Set(headerKey, headerValue)

	return crlfPos + 2, false, nil
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h *Headers) Set(key, value string) {
	name := strings.ToLower(key)

	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) ForEach(fn func(key, value string)) {
	for key, value := range h.headers {
		fn(key, value)
	}
}

var validBytePattern = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_\` + "`" + `|~]*$`)

// IsValidByteSlice checks if the given byte slice contains only allowed characters:
// - Uppercase letters: A-Z
// - Lowercase letters: a-z
// - Digits: 0-9
// - Special characters: !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~
func (h *Headers) isValidByteSlice(data []byte) bool {
	return validBytePattern.Match(data)
}
