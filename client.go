package gearman

import (
	"sync"
)

type client struct {
	mutex     sync.Mutex
	hosts     []string
	hostState []hostState
	
	handlers  map[string]func(job *IncomingJob) []byte
}

type hostState struct {
}

