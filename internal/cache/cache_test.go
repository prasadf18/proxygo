package cache

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := NewCache(1 * time.Minute)

	c.Set("key1", []byte("hello"))

	data, found := c.Get("key1")
	if !found {
		t.Fatal("expected to find key1, but it was not found")
	}
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", string(data))
	}
}

func TestGetExpiredEntry(t *testing.T) {
	c := NewCache(10 * time.Millisecond)

	c.Set("key1", []byte("hello"))

	time.Sleep(20 * time.Millisecond)

	_, found := c.Get("key1")
	if found {
		t.Error("expected key1 to be expired, but it was still found")
	}
}

func TestGetMissingKey(t *testing.T) {
	c := NewCache(1 * time.Minute)

	_, found := c.Get("does-not-exist")
	if found {
		t.Error("expected key to not be found, but it was")
	}
}