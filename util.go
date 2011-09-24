package gearman
import (
	"os"
	"net"
	"io"
	"log"
)

const (
        CAN_DO = 1                       // =>  [ 'I', "can_do" ],     # from W:  [FUNC]
        CAN_DO_TIMEOUT = 23              // => [ 'I', "can_do_timeout" ], # from W: FUNC[0]TIMEOUT
        CANT_DO = 2                      // =>  [ 'I', "cant_do" ],    # from W:  [FUNC]
        RESET_ABILITIES = 3              // =>  [ 'I', "reset_abilities" ],  # from W:  ---
        SET_CLIENT_ID = 22               // => [ 'I', "set_client_id" ],    # W->J: [RANDOM_STRING_NO_WHITESPACE]
        PRE_SLEEP = 4                    // =>  [ 'I', "pre_sleep" ],  # from W: ---

        NOOP = 6                         // =>  [ 'O', "noop" ],        # J->W  ---
        SUBMIT_JOB = 7                   // =>  [ 'I', "submit_job" ],    # C->J  FUNC[0]UNIQ[0]ARGS
        SUBMIT_JOB_HIGH = 21             // =>  [ 'I', "submit_job_high" ],    # C->J  FUNC[0]UNIQ[0]ARGS
        SUBMIT_JOB_BG = 18               // => [ 'I', "submit_job_bg" ], # C->J     " "   "  " "

        JOB_CREATED = 8                  // =>  [ 'O', "job_created" ], # J->C HANDLE
        GRAB_JOB = 9                     // =>  [ 'I', "grab_job" ],    # W->J --
        NO_JOB = 10                      // => [ 'O', "no_job" ],      # J->W --
        JOB_ASSIGN = 11                  // => [ 'O', "job_assign" ],  # J->W HANDLE[0]FUNC[0]ARG

        WORK_STATUS = 12                 // => [ 'IO',  "work_status" ],   # W->J/C: HANDLE[0]NUMERATOR[0]DENOMINATOR
        WORK_COMPLETE = 13               // => [ 'IO',  "work_complete" ], # W->J/C: HANDLE[0]RES
        WORK_FAIL = 14                   // => [ 'IO',  "work_fail" ],     # W->J/C: HANDLE

        GET_STATUS = 15                  // => [ 'I',  "get_status" ],  # C->J: HANDLE
        STATUS_RES = 20                  // => [ 'O',  "status_res" ],  # C->J: HANDLE[0]KNOWN[0]RUNNING[0]NUM[0]DENOM

        ECHO_REQ = 16                    // => [ 'I',  "echo_req" ],    # ?->J TEXT
        ECHO_RES = 17                    // => [ 'O',  "echo_res" ],    # J->? TEXT

        ERROR = 19                       // => [ 'O',  "error" ],       # J->? ERRCODE[0]ERR_TEXT

	//            # for worker to declare to the jobserver that this worker is only connected
	//            # to one jobserver, so no polls/grabs will take place, and server is free
	//            # to push "job_assign" packets back down.
	ALL_YOURS = 24                    // => [ 'I', "all_yours" ],    # W->J ---
)

func make_req ( cmd int, d []byte ) []byte {
	buf := make ( []byte, 0 )
	buf = append ( buf, 0 )
	buf = append ( buf, []byte ( "REQ" ) ... )
	buf = append ( buf, uint8(cmd >> 24 ), uint8(cmd>>16 ),
                uint8 ( cmd>>8  ), uint8 ( cmd ) )
	l := len ( d )
	buf = append  (buf, uint8(l>>24), uint8(l>>16), uint8(l>>8), uint8(l))
	return append ( buf, d ... )
}

func make_res ( cmd int, d []byte ) []byte {
	buf := make ( []byte, 0 )
	buf = append ( buf, 0 )
	buf = append ( buf, []byte ( "RES" ) ... )
	buf = append ( buf, uint8(cmd >> 24 ), uint8(cmd>>16 ),
                uint8 ( cmd>>8  ), uint8 ( cmd ) )
	l := len ( d )
	buf = append  (buf, uint8(l>>24), uint8(l>>16), uint8(l>>8), uint8(l))
	return append ( buf, d ... )
}

func (c *client) client_id () []byte {
	c.id = "lczhtjrfiwwengeyyplvuxtkubfiwa"
	return make_req ( SET_CLIENT_ID, []byte ( c.id ) )
}

func read_header ( n net.Conn ) ( uint32, uint32, bool, os.Error ) {
	headerbuf := make ( []byte, 12 )
	_, e := io.ReadFull ( n, headerbuf )
	if e != nil {
		if to := e.(net.Error); to != nil && to.Timeout() {
			if e!= nil {
				return 0, 0, true, nil
			}
			return 0, 0, false, e
		}
	}
	log.Println ( headerbuf )
	cmd := (uint32(headerbuf[4]) << 24) | (uint32(headerbuf[5]) << 16) |
		(uint32(headerbuf[6]) << 8) | uint32(headerbuf[7])
	cmd_len := (uint32(headerbuf[8]) << 24) | (uint32(headerbuf[9]) << 16) |
		(uint32(headerbuf[10]) << 8) | uint32(headerbuf[11])
	return cmd, cmd_len, false, nil
}
