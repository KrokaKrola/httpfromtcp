package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp.krokakrola.com/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("error listening for TCP: %s\n", err.Error())
	}

	fmt.Printf("Started TCP connection on port: %s\n", port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}

		fmt.Println("Accepted connection from:", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatalf("error reading request from reader: %s\n", err.Error())
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		req.Headers.ForEach(func(key, value string) {
			fmt.Printf("- %s: %s\n", key, value)
		})
		fmt.Println("Body:")
		fmt.Printf("%s\n", string(req.Body))

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
