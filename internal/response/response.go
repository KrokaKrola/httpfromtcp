package response

import (
	"fmt"
	"io"

	"httpfromtcp.krokakrola.com/internal/headers"
)

const protocol = "HTTP"
const protocolVersion = "1.1"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("%s/%s %d %s\r\n", protocol, protocolVersion, statusCode, statusCode.String())

	_, err := w.Write([]byte(statusLine))

	return err
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	headers.ForEach(func(key, value string) {
		w.Write(fmt.Appendf(nil, "%s: %s\r\n", key, value))
	})

	_, err := w.Write([]byte("\r\n"))

	return err
}
