package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jgrove2/tcp-servers/internal/request"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}
		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("Request line\nMethod: %s\nTarget: %s\nVersion: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range request.Headers {
			fmt.Printf("%s: %s\n", k, v)
		}
		fmt.Printf("Body: \n%s", request.Body)
	}

}
