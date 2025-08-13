package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)

			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				ch <- str
				str = ""
			}

			str += string(data)
		}

		if len(str) != 0 {
			ch <- str
		}
	}()

	return ch
}

const port = 42069
const host = "127.0.0.1"

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Started tcp connection on port: %d", port)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection accepted")
		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Printf("read: %s\n", line)
		}

		fmt.Println("Connection closed")
	}
}
