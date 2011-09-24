package gearman

import (
	"sync"
	"net"
)

type client struct {
	mutex     sync.Mutex
	hosts     []string
	hostState []hostState
	
	handlers  map[string]func(job *IncomingJob) []byte
	id string
}

type hostState struct {
	conn net.Conn
}

