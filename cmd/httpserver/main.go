package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jgrove2/tcp-servers/internal/server"
)

const port = 8080

func main() {
	httpServer := server.NewServer()
	httpServer.Serve(port)
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
	defer httpServer.Close()
}
