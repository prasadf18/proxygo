package proxy

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleConnectionForwardsRequest(t *testing.T) {
	// Fake upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from fake upstream"))
	}))
	defer upstream.Close()

	upstreamAddr := upstream.Listener.Addr().String()

	// Start a listener for the proxy on a random free port
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
			go HandleConnection(conn, upstreamAddr)
		}
	}()

	time.Sleep(100 * time.Millisecond) // give the listener a moment

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