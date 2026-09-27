package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"net/http"
)

func main() {
	target := flag.String("target", "localhost:9000", "upstream address to proxy to")
	flag.Parse()

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	log.Printf("listening on :8080, forwarding to %s", *target)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConnection(conn, *target)
	}
}

func handleConnection(conn net.Conn, target string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Printf("failed to read request: %v", err)
		return
	}

	log.Printf("received request: %s %s (target: %s)", req.Method, req.URL.Path, target)

	body := "hello from proxy"
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 16\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		body

	conn.Write([]byte(response))
}