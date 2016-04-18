package xerrors

/*
* Work in progress (custom) error handling that aims to generalize the error types for easier error handling 
*/

import (
	"runtime"
	"strconv"
	)


var sysErrMsg = "Application Error!"
var peerErrMsg = "Peer Error!"
var sysErrCode = `001`
var peerErrCode = `002`

/*
* The system errors are FATAL, UNRECOVERABLE errors in the BUSINESS logic of the application. The system error should be considered as an ASSERT.
* The CURRENT executed job should be aborted.
* When a system error appears it means that there is a BUG in the code. 
* The command that caused the system error can be retried.
*/

func NewSysErr() XError{
    _, fn, line, _ := runtime.Caller(1)
    return XError{sysErrMsg + ` file name: ` + fn + ` line: ` + strconv.Itoa(line), sysErrMsg, sysErrCode, true}
}

/*
* The peer errors are FATAL, UNRECOVERABLE errors in a REMOTE system.
* The CURRENT executed job should be aborted.
* When a Peer error apears it means that the remote system is down.
* A general message should be displayed to the end user with the name (or unique identifier) of the remote system
* The command that caused the peer error can be retried.
*/

func NewPeerErr(debugMsg string) XError{
    return XError{debugMsg, peerErrMsg, peerErrCode, true}
}


type XError struct{
        UIMsg		string // Message that should be displayed to the user, saying what the user should do
	DebugMsg	string // Message for the developer
	Code		string
	IsRetryable	bool
}

func (e XError) Error() string {
        return  e.UIMsg
}


/*
* The UI errors may be FATAL, but can also be RECOVERABLE. It is should be considered as an EXCEPTION. 
* If the current executed job is aborted depends on the error itself.
*/
func NewUIErr(uiMsg string, debugMsg string, code string, IsRetryable bool) XError{
        return XError{uiMsg, debugMsg, code, IsRetryable}
}


