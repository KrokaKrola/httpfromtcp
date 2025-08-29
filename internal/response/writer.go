package response

import (
	"fmt"
	"io"

	"httpfromtcp.krokakrola.com/internal/headers"
)

type Writer struct {
	conn io.Writer
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		conn,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := fmt.Sprintf("%s/%s %d %s\r\n", "HTTP", "1.1", statusCode, statusCode.String())

	_, err := w.conn.Write([]byte(statusLine))

	return err
}

func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	headers.ForEach(func(key, value string) {
		w.conn.Write(fmt.Appendf(nil, "%s: %s\r\n", key, value))
	})

	_, err := w.conn.Write([]byte("\r\n"))

	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteBody(body []byte) error {
	_, err := w.conn.Write(body)

	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteHtml(statusCode StatusCode, body string) error {
	htmlStr := fmt.Sprintf(
		"<html><head><title>%d %s</title></head><body><h1>%s</h1><p>%s</p></body></html>",
		statusCode, statusCode.String(), statusCode.String(), body,
	)

	headers := headers.NewHeaders()

	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/html")
	headers.Set("Content-Length", fmt.Sprint(len(htmlStr)))

	err := w.WriteStatusLine(statusCode)

	if err != nil {
		return err
	}

	err = w.WriteHeaders(headers)

	if err != nil {
		return err
	}

	err = w.WriteBody([]byte(htmlStr))

	return err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.conn, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.conn.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.conn.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.conn.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *Writer) WriteTrailers(h *headers.Headers) error {
	h.ForEach(func(key, value string) {
		fmt.Fprintf(w.conn, "%s: %s\r\n", key, value)
	})

	_, err := fmt.Fprintf(w.conn, "\r\n")

	if err != nil {
		return err
	}

	return nil
}
