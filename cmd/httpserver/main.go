package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jgrove2/tcp-servers/internal/request"
	"github.com/jgrove2/tcp-servers/internal/server"
)

const port = 8080

func handler(req *request.Request, w io.Writer) *server.HandleError {
	log.Println(req)
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandleError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandleError{
			StatusCode: 500,
			Message:    "Woopsie, my bad\n",
		}
	default:
		fmt.Fprint(w, "Hello World!")
		break
	}
	return nil
}

func main() {
	httpServer := server.NewServer()
	httpServer.Serve(8080, handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
	defer httpServer.Close()
}
