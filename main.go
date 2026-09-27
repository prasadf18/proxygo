package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	log.Println("listening on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Printf("failed to read request: %v", err)
		return
	}

	log.Printf("received request: %s %s", req.Method, req.URL.Path)

	body := "hello from proxy"
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 17\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		body

	conn.Write([]byte(response))
}