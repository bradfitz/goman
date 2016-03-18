package gearman

import (
	"net"
	"testing"
)

const TEST_ROUTER = "localhost:4730"

func routerIsUp(hostPort string) bool {
	c, err := net.Dial("tcp", "", hostPort)
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func TestFoo(t *testing.T) {
	if !routerIsUp(TEST_ROUTER) {
		t.Fatalf("Can't run unit tests without gearmand running on %s", TEST_ROUTER)
	}

	c := NewClient(TEST_ROUTER)
	if c == nil {
		t.Fatal("Got nil client")
	}
}
