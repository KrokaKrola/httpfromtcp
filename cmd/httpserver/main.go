package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"httpfromtcp.krokakrola.com/internal/request"
	"httpfromtcp.krokakrola.com/internal/response"
	"httpfromtcp.krokakrola.com/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(req *request.Request, res *response.Writer) {
		if req.RequestLine.RequestTarget == "/" {
			res.WriteHtml(response.OK, "Your request was an absolute banger.")
		}

		if req.RequestLine.RequestTarget == "/yourproblem" {
			res.WriteHtml(response.BadRequest, "Your request honestly kinda sucked.")
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			res.WriteHtml(response.InternalServerError, "Okay, you know what? This one is on me.")
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
