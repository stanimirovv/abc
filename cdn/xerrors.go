package main

/*
* This file contains the custom errors and logging that will be used in the project.
* It will eventually become a package, when things have been figured out, until then it will be a "plug and play" file.
*/



/*
* The system errors are FATAL, UNRECOVERABLE errors in the BUSINESS logic of the application.
* The current execution path must be aborted. 
* If the process ends or is freed to start working on the next job is up to the developer.
*/

type SysErr struct {
    msg string
}

func (e *SysErr) Error() string {
    //return fmt.Sprintf(" %s", e.prob)
    return  "Application Error!"
}

func NewSysErr(msg string) SysErr{
    return SysErr{msg}
}

