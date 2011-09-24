package gearman

import (
	"testing"
	"net"
)

const TEST_ROUTER = "localhost:6481"

func routerIsUp(hostPort string) bool {
	c, err := net.Dial("tcp", hostPort)
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

	go func () {
		w := NewClient(TEST_ROUTER)
		if w == nil {
			t.Fatal("Got nil client")
		}
		w.RegisterWorker ( "geturl", geturl )
		w.Work()
	}()
	servers := []string{TEST_ROUTER}
	c := NewLoadBalancedClient(servers)
	if c == nil {
		t.Fatal("Got nil client")
	}
	res := c.Call ( "geturl", []byte ( "http://tinychat.com" ) )
	if res == nil {
		t.Fatal ( "Bad response" )
	}
}

func geturl ( job *IncomingJob ) []byte {
	url := string ( job.Data )
	return []byte ( url )
}