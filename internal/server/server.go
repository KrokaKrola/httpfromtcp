package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"httpfromtcp.krokakrola.com/internal/request"
	"httpfromtcp.krokakrola.com/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			if s.closed.Load() {
				return
			}

			log.Printf("Error accepting connection: %+v\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from:", conn.RemoteAddr())

	res := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		res.WriteHtml(response.BadRequest, "Your request honestly kinda sucked.")
		return
	}

	s.handler(req, res)
}
