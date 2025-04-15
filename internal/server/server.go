package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jgrove2/tcp-servers/internal/request"
	"github.com/jgrove2/tcp-servers/internal/response"
)

type IServer interface {
	Close() error
	Serve(port int) error
	Listen()
}

type Server struct {
	closed   bool
	listener net.Listener
	handler  HandlerFunc
}

type HandleError struct {
	StatusCode int
	Message    string
}

type HandlerFunc func(req *request.Request, w io.Writer) *HandleError

func NewServer() *Server {
	return &Server{
		closed:   false,
		listener: nil,
		handler:  nil,
	}
}

func (s *Server) Close() error {
	s.closed = true
	return s.listener.Close()
}

func (s *Server) Serve(port int, handler HandlerFunc) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	s.handler = handler
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

func writeHandlerError(conn io.Writer, herr *HandleError) {
	response.WriteStatusLine(conn, response.StatusCode(herr.StatusCode))
	headers := response.GetDefaultHeaders(len(herr.Message))
	response.WriteHeaders(conn, headers)
	conn.Write([]byte(herr.Message))
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer log.Println("Closing connection")

	req, err := request.RequestFromReader(conn)
	if err != nil {
		writeHandlerError(conn, &HandleError{
			StatusCode: 400,
			Message:    "Bad Request",
		})
		return
	}

	var buf bytes.Buffer
	log.Println(req)
	if herr := s.handler(req, &buf); herr != nil {
		writeHandlerError(conn, herr)
		return
	}

	body := buf.Bytes()
	hders := response.GetDefaultHeaders(len(body))
	err = response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		return
	}
	err = response.WriteHeaders(conn, hders)
	if err != nil {
		return
	}
	conn.Write(body)
}
