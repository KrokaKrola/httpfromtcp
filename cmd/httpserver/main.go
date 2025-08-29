package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"httpfromtcp.krokakrola.com/internal/headers"
	"httpfromtcp.krokakrola.com/internal/request"
	"httpfromtcp.krokakrola.com/internal/response"
	"httpfromtcp.krokakrola.com/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(req *request.Request, res *response.Writer) {

		if req.RequestLine.RequestTarget == "/yourproblem" {
			res.WriteHtml(response.BadRequest, "Your request honestly kinda sucked.")
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			res.WriteHtml(response.InternalServerError, "Okay, you know what? This one is on me.")
		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
			r, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", target))

			if err != nil {
				res.WriteHtml(response.InternalServerError, "httpbin error")
				return
			}

			defer r.Body.Close()

			res.WriteStatusLine(response.OK)

			h := response.GetDefaultHeaders(0)

			h.Delete("content-length")
			h.Replace("Transfer-Encoding", "chunked")
			h.Replace("Trailer", "X-Content-SHA256")
			h.Set("Trailer", "X-Content-Length")

			res.WriteHeaders(h)

			b := make([]byte, 1024)
			body := make([]byte, 0)
			for {
				n, err := r.Body.Read(b)
				fmt.Println("Read", n, "bytes")

				if n > 0 {
					_, err = res.WriteChunkedBody(b[:n])

					if err != nil {
						log.Println("Error while reading response body:", err)
						break
					}

					body = append(body, b[:n]...)
				}

				if err == io.EOF {
					break
				}

				if err != nil {
					log.Println("Error while reading response body:", err)
				}
			}
			_, err = res.WriteChunkedBodyDone()
			if err != nil {
				log.Println("error writing chuncked body done", err)
			}

			trailers := headers.NewHeaders()
			sum := sha256.Sum256(body)
			trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", sum))
			trailers.Set("X-Content-Length", fmt.Sprint(len(body)))
			err = res.WriteTrailers(trailers)
			if err != nil {
				log.Println("error writing trailers", err)
			}
		} else if req.RequestLine.RequestTarget == "/video" {
			h := response.GetDefaultHeaders(0)
			h.Replace("Content-Type", "video/mp4")
			h.Delete("Connection")
			b, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				res.WriteHtml(response.InternalServerError, "File was not found")
				return
			}
			h.Replace("Content-Length", fmt.Sprint(len(b)))
			log.Printf("headers: %+v", h)
			res.WriteStatusLine(response.OK)
			res.WriteHeaders(h)
			res.WriteBody(b)
		} else {
			res.WriteHtml(response.OK, "Your request was an absolute banger.")
		}
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
