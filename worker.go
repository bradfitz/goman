// The "worker" side of a Gearman client.

package gearman

import (
	"log"
)

func (j *IncomingJob) SetProgress(done int, total int) {
	// TODO: send progress packets
}

func (c *client) RegisterWorker(method string, handler func(job *IncomingJob) []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.handlers == nil {
		c.handlers = make(map[string]func(job *IncomingJob) []byte)
	}
	c.handlers[method] = handler
}

func (c *client) Work() {
	log.Exitf("TODO: implement Work()")

	// TODO: send CANDO packets to job servers
	// TODO: block in a loop waiting on work packets
}
