package gearman

import (
	"testing"
	"net"
	"http"
	"io/ioutil"
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

	go func() {
		w := NewClient(TEST_ROUTER)
		if w == nil {
			t.Fatal("Got nil client")
		}
		w.RegisterWorker("geturl", geturl)
		w.Work()
	}()
	servers := []string{TEST_ROUTER}
	c := NewLoadBalancedClient(servers)
	if c == nil {
		t.Fatal("Got nil client")
	}
	res := c.Call("geturl", []byte("http://tinychat.com/c0afc7f20626fbbdd705d2e1bd7bc5dd370dbdd0"))
	if len(res) == 0 {
		t.Fatal("No response")
	}
	if string(res) != "gearman was here\n" {
		t.Fatal("Bad response")
	}
}

func geturl(job *IncomingJob) []byte {
	url := string(job.Data)
	resp, e := http.Get(url)
	if e != nil {
		return nil
	}
	body, e := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		return nil
	}
	return body
}
