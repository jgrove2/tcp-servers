package server

import (
	"fmt"
	"log"
	"net"

	"github.com/jgrove2/tcp-servers/internal/request"
)

type IServer interface {
	Close() error
	Serve(port int) error
	Listen()
}

type Server struct {
	closed   bool
	listener net.Listener
}

func NewServer() *Server {
	return &Server{
		closed:   false,
		listener: nil,
	}
}

func (s *Server) Close() error {
	s.closed = true
	return s.listener.Close()
}

func (s *Server) Serve(port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	go s.Listen()
	return nil
}

func (s *Server) Listen() {
	for {
		if s.closed {
			break
		}

		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed {
				return
			}
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer log.Println("Closing connection")

	log.Println("New connection from", conn.RemoteAddr())
	req, err := request.RequestFromReader(conn)
	log.Println(req)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 11\r\n\r\nHello world"))
	if err != nil {
		log.Println(err.Error())
		return
	}
}
