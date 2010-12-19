// The "caller" side of a Gearman client.

package gearman

import (
	"log"
)

func (c *client) Call(method string, data []byte) []byte {
	return c.CallWithProgress(method, data, nil)
}

func (c *client) CallWithProgress(method string, data []byte, progress ProgressHandler) []byte {
	// TODO: implement

	// TODO: shard onto a c.hosts[] entry, wait for its connection state, etc.
	log.Exitf("TODO: implement CallWithProgress")
	return nil
}
