package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/prasadf18/proxygo/internal/cache"
	"github.com/prasadf18/proxygo/internal/proxy"
)

func main() {
	target := flag.String("target", "localhost:9000", "upstream address to proxy to")
	flag.Parse()

	c := cache.NewCache(30 * time.Second)

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

		go proxy.HandleConnection(conn, *target, c)
	}
}