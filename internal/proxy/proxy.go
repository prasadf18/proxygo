package proxy

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"net/http"

	"github.com/prasadf18/proxygo/internal/cache"
)

func HandleConnection(conn net.Conn, target string, c *cache.Cache) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Printf("failed to read request: %v", err)
		return
	}

	log.Printf("received request: %s %s (target: %s)", req.Method, req.URL.Path, target)

	key := req.Method + " " + req.URL.String()

	if req.Method == http.MethodGet {
		if cached, found := c.Get(key); found {
			log.Printf("cache HIT for %s", key)
			_, err = conn.Write(cached)
			if err != nil {
				log.Printf("failed to write cached response: %v", err)
			}
			return
		}

		log.Printf("cache MISS for %s", key)
	}

	upstreamConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("failed to connect to upstream: %v", err)
		return
	}
	defer upstreamConn.Close()

	err = req.Write(upstreamConn)
	if err != nil {
		log.Printf("failed to forward request: %v", err)
		return
	}

	upstreamReader := bufio.NewReader(upstreamConn)
	resp, err := http.ReadResponse(upstreamReader, req)
	if err != nil {
		log.Printf("failed to read upstream response: %v", err)
		return
	}

	var buf bytes.Buffer

	err = resp.Write(&buf)
	if err != nil {
		log.Printf("failed to serialize response: %v", err)
		return
	}

	if req.Method == http.MethodGet && resp.StatusCode == http.StatusOK {
		c.Set(key, buf.Bytes())
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("failed to relay response to client: %v", err)
		return
	}
}