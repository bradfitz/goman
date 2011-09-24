// The "caller" side of a Gearman client.

//        SUBMIT_JOB = 7                   // =>  [ 'I', "submit_job" ],    # C->J  FUNC[0]UNIQ[0]ARGS
//        SUBMIT_JOB_HIGH = 21             // =>  [ 'I', "submit_job_high" ],    # C->J  FUNC[0]UNIQ[0]ARGS
//        SUBMIT_JOB_BG = 18               // => [ 'I', "submit_job_bg" ], # C->J     " "   "  " "

//        JOB_CREATED = 8                  // =>  [ 'O', "job_created" ], # J->C HANDLE

//        WORK_STATUS = 12                 // => [ 'IO',  "work_status" ],   # W->J/C: HANDLE[0]NUMERATOR[0]DENOMINATOR
//        WORK_COMPLETE = 13               // => [ 'IO',  "work_complete" ], # W->J/C: HANDLE[0]RES
//        WORK_FAIL = 14                   // => [ 'IO',  "work_fail" ],     # W->J/C: HANDLE

//        GET_STATUS = 15                  // => [ 'I',  "get_status" ],  # C->J: HANDLE
//        STATUS_RES = 20                  // => [ 'O',  "status_res" ],  # C->J: HANDLE[0]KNOWN[0]RUNNING[0]NUM[0]DENOM

//        ECHO_REQ = 16                    // => [ 'I',  "echo_req" ],    # ?->J TEXT
//        ECHO_RES = 17                    // => [ 'O',  "echo_res" ],    # J->? TEXT

//        ERROR = 19                       // => [ 'O',  "error" ],       # J->? ERRCODE[0]ERR_T
package gearman

import (
	"net"
	"io"
	"os"
	"time"
	"bytes"
	"rand"
)

func (c *client) Call(method string, data []byte) []byte {
	return c.CallWithProgress(method, data, nil)
}

func (c *client) CallWithProgress(method string, data []byte, progress ProgressHandler) []byte {
	maxtries := len(c.hosts) * 2
	var n net.Conn = nil
	var jobhandle []byte = nil
	rand.Seed(time.Nanoseconds())

	// find a jobserver that will handle this method
	for maxtries > 0 {
		maxtries--
		rnum := rand.Intn(len(c.hosts))
		// is this conn alive?
		if !(c.hostState[rnum].conn != nil && c.hostState[rnum].conn.RemoteAddr() != nil) {
			var e os.Error = nil
			n, e = net.Dial("tcp", c.hosts[rnum])
			if e != nil {
				//log.Println ( "finding a job server " + e.String() )
				continue
			}
		}
		c.hostState[rnum].conn = n
		buf := []byte(method)
		buf = append(buf, 0)
		c.id = "jfidfjid"
		buf = append(buf, []byte(c.id)...)
		buf = append(buf, 0)
		buf = append(buf, data...)
		_, e := n.Write(make_req(SUBMIT_JOB, buf))
		if e != nil {
			n.Close()
			continue
		}
		cmd, cmd_len, to, e := read_header(n)
		if e != nil || to {
			n.Close()
			continue
		}
		data := make([]byte, cmd_len)
		_, e = io.ReadFull(n, data)
		if e != nil {
			n.Close()
			continue
		}
		if cmd != JOB_CREATED {
			continue
		}
		jobhandle = data
		break
	}
	if jobhandle == nil {
		return nil
	}

	for {
		cmd, cmd_len, to, e := read_header(n)
		data := make([]byte, cmd_len)
		_, e = io.ReadFull(n, data)
		if e != nil {
			return nil
		}
		if to {
			continue
		}
		switch cmd {
		case WORK_COMPLETE:
			if len(data) == 0 {
				return nil
			}
			a := bytes.SplitN(data, []byte{0}, 2)
			if len(a) != 2 {
				return nil
			}
			return a[1]
		}
	}
	return nil
}
