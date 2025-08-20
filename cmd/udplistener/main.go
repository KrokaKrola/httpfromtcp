package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Remote address (where to send messages)
	remoteAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to remote UDP server (nc listener)
	// Pass nil as local address to let Go choose an available port
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Send the message to the remote UDP server
		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
