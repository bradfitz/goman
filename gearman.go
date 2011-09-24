package gearman

type Job struct {
	Method string
	Data   []byte
}

type IncomingJob struct {
	*Job
}

type ProgressHandler interface {
	OnProgress(done int, total int)
}

type Client interface {
	// For being a worker:
	RegisterWorker(method string, handler func(job *IncomingJob) []byte)
	Work()

	// For being a client:
	Call(method string, data []byte) []byte
	CallWithProgress(method string, data []byte, progress ProgressHandler) []byte
}

func (ij *IncomingJob) SendProgress(done int, total int) {
	// TODO: implement
}

func NewClient(hostport string) Client {
	return &client{hosts: []string{hostport}, hostState: make ( []hostState, 1 )}
}

func NewLoadBalancedClient(hostports []string) Client {
	return &client{hosts: hostports, hostState: make ( []hostState, len ( hostports ) ) }
}

