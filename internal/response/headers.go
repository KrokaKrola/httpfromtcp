package response

import (
	"fmt"

	"httpfromtcp.krokakrola.com/internal/headers"
)

const defaultContentType = "text/plain"

func GetDefaultHeaders(contentLen int) *headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Connection", "close")
	headers.Set("Content-Type", defaultContentType)
	headers.Set("Content-Length", fmt.Sprint(contentLen))

	return headers
}
