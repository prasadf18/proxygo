package proxy

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prasadf18/proxygo/internal/cache"
)

func TestHandleConnectionForwardsRequest(t *testing.T) {
	
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from fake upstream"))
	}))
	defer upstream.Close()

	upstreamAddr := upstream.Listener.Addr().String()

	c := cache.NewCache(30 * time.Second)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start proxy listener: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go HandleConnection(conn, upstreamAddr, c)
		}
	}()

	time.Sleep(100 * time.Millisecond) 

	resp, err := http.Get("http://" + ln.Addr().String() + "/test")
	if err != nil {
		t.Fatalf("request to proxy failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(body) != "hello from fake upstream" {
		t.Errorf("expected 'hello from fake upstream', got %q", string(body))
	}
}