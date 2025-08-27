package response

import (
	"fmt"
	"io"

	"httpfromtcp.krokakrola.com/internal/headers"
)

const protocol = "HTTP"
const protocolVersion = "1.1"
const defaultContentType = "text/plain"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("%s/%s %d %s\r\n", protocol, protocolVersion, statusCode.Number(), statusCode.String())

	_, err := w.Write([]byte(statusLine))

	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Connection", "close")
	headers.Set("Content-Type", defaultContentType)
	headers.Set("Content-Length", fmt.Sprint(contentLen))

	return headers
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	err := headers.ForEach(func(key, value string) error {
		_, err := w.Write(fmt.Appendf(nil, "%s: %s\r\n", key, value))

		return err
	})

	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\r\n"))

	if err != nil {
		return err
	}

	return nil
}
