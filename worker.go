// The "worker" side of a Gearman client.

// from Gearman::Util.pm
//        CAN_DO = 1                       // =>  [ 'I', "can_do" ],     # from W:  [FUNC]
//        CAN_DO_TIMEOUT = 23              // => [ 'I', "can_do_timeout" ], # from W: FUNC[0]TIMEOUT
//        CANT_DO = 2                      // =>  [ 'I', "cant_do" ],    # from W:  [FUNC]
//        RESET_ABILITIES = 3              // =>  [ 'I', "reset_abilities" ],  # from W:  ---
//        SET_CLIENT_ID = 22               // => [ 'I', "set_client_id" ],    # W->J: [RANDOM_STRING_NO_WHITESPACE]
//        PRE_SLEEP = 4                    // =>  [ 'I', "pre_sleep" ],  # from W: ---
//        GRAB_JOB = 9                     // =>  [ 'I', "grab_job" ],    # W->J --
//        WORK_STATUS = 12                 // => [ 'IO',  "work_status" ],   # W->J/C: HANDLE[0]NUMERATOR[0]DENOMINATOR
//        WORK_COMPLETE = 13               // => [ 'IO',  "work_complete" ], # W->J/C: HANDLE[0]RES
//        WORK_FAIL = 14                   // => [ 'IO',  "work_fail" ],     # W->J/C: HANDLE
//        ECHO_REQ = 16                    // => [ 'I',  "echo_req" ],    # ?->J TEXT
//	ALL_YOURS = 24                    // => [ 'I', "all_yours" ],    # W->J ---


//        NO_JOB = 10                      // => [ 'O', "no_job" ],      # J->W --
//        NOOP = 6                         // =>  [ 'O', "noop" ],        # J->W  ---
//        JOB_ASSIGN = 11                  // => [ 'O', "job_assign" ],  # J->W HANDLE[0]FUNC[0]ARG
//        ECHO_RES = 17                    // => [ 'O',  "echo_res" ],    # J->? TEXT
//        ERROR = 19                       // => [ 'O',  "error" ],       # J->? ERRCODE[0]ERR_TEXT



package gearman

import (
	"net"
	"io"
	"time"
	"os"
	"bytes"
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
	c.id = "jfidjfid"
	for _, v := range c.hosts {
		n, e := net.Dial("tcp", v)
		if e != nil {
			continue
		}
		_, e = n.Write(c.client_id())
		if e != nil {
			return
		}
		go c.worker_loop(n)
	}
	for {
		time.Sleep(1e9 * 5)
	}
}

func (c *client) worker_loop(n net.Conn) {
	defer n.Close()
	for h, _ := range c.handlers {
		buf := make([]byte, 0)
		buf = append(buf, []byte(h)...)
		buf = append(buf, 0)
		_, e := n.Write(make_req(CAN_DO, []byte(h)))
		if e != nil {
			return
		}
	}
	n.SetReadTimeout(1e9 * 60)
	_, e := n.Write(make_req(GRAB_JOB, []byte{}))
	if e != nil {
		return
	}
	for {
		// worker asks for jobs periodically
		// server will only push a job to a sleeping worker
		cmd, cmd_len, to, e := read_header(n)
		if e != nil {
			return
		}
		if to {
			// timed out, ask for another job
			_, e := n.Write(make_req(GRAB_JOB, []byte{}))
			if e != nil {
				return
			}
			continue
		}
		if cmd == NO_JOB {
			// no jobs, go asleep
			_, e = n.Write(make_req(PRE_SLEEP, []byte{}))
			if e != nil {
				return
			}
			continue
		}

		databuf := make([]byte, cmd_len)
		_, e = io.ReadFull(n, databuf)
		if e != nil {
			return
		}

		switch cmd {
		case NOOP:
			// a wakeup call?
		case ECHO_RES:
		case ERROR:
			//a := bytes.SplitN ( databuf, []byte{0}, 2 )
		case JOB_ASSIGN:
			done := false
			var buf []byte = nil
			for !done {
				buf, done, e = c.do_work(cmd, databuf)
				if e != nil {
					break
				}
			}
			if buf == nil {
				buf = []byte{1, 0}
			}
			_, e = n.Write(make_req(WORK_COMPLETE, buf))
			if e != nil {
				return
			}
		}
		_, e = n.Write(make_req(GRAB_JOB, []byte{}))
		if e != nil {
			return
		}
	}
}

func (c *client) do_work(cmd uint32, data []byte) ([]byte, bool, os.Error) {
	a := bytes.SplitN(data, []byte{0}, 3)
	if len(a) != 3 {
		return []byte{1, 0}, true, os.NewError("not enough args")
	}
	buf := a[0] // handle
	buf = append(buf, 0)
	f, ok := c.handlers[string(a[1])]
	if !ok {
		return buf, true, os.NewError("this worker does not handle " + string(a[1]))
	}
	res := f(&IncomingJob{&Job{string(a[1]), a[2]}})
	if res != nil && len(res) > 0 {
		buf = append(buf, res...)
		return buf, true, nil
	}
	return buf, true, os.NewError("dont know")
}
