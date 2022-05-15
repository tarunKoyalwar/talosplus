package gopool

import "github.com/tarunKoyalwar/talosplus/pkg/shell"

// JobResponse :  Containing Response of A Executed Function
type JobResponse struct {
	Err      error
	Uid      string  //Task UID
	ExecTime float32 //Execution Time
}

// Job : Defining the Job
type Job struct {
	Cx  *shell.CMDWrap //  pointer to command struct
	UID string
}
