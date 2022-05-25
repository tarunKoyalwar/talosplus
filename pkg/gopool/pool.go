package gopool

import (
	"sync"

	"github.com/tarunKoyalwar/talosplus/pkg/shell"
)

// Pool : Goroutine Pool with Limit
type Pool struct {
	Concurrency int
	JobCount    int

	workerChannel chan int         // Used to Limit Number of Active Goroutines
	JobChannel    chan Job         // Channel Containing Functions to be executed
	InterChannel  chan JobResponse // Channel Between Launcher and Operator
	Completion    chan int         // Will Block until all jobs are completed

	HandleError  func(er error)         //What to do if Job Execution Failed
	OnCompletion func(resp JobResponse) //Any Furthur steps after job completion

	logic    *sync.WaitGroup // Handles shared Working of Pool
	Wg       *sync.WaitGroup // Wait Group For Launched Jobs
	JobMutex *sync.Mutex
	JobSent  int
}

// AddJob
func (p *Pool) AddJob(t *shell.CMDWrap) {
	j := Job{
		Cx:  t,
		UID: "NA",
	}

	p.JobChannel <- j
	p.JobSent += 1
}

// AddJobWithId
func (p *Pool) AddJobWithId(t *shell.CMDWrap, uid string) {
	j := Job{
		Cx:  t,
		UID: uid,
	}

	p.JobChannel <- j
	p.JobSent += 1
}

// Launcher : (shared)Launches New Jobs assings workers etc
func (p *Pool) Launcher() {
	defer p.logic.Done()

	//create worker based on concurrency
	for i := 0; i < p.Concurrency; i++ {
		p.workerChannel <- 8 //random number
	}

	for {
		rjob, ok := <-p.JobChannel

		//If Channel Is Closed then break away
		if !ok {
			break
		}

		//check if Any Worker is Available
		//Will be blocked if No Workers are available
		_, status := <-p.workerChannel
		if !status {
			//workers quit so break away
			break
		}

		// fmt.Println("Starting New JOb")
		//Complete the Work
		p.Wg.Add(1)
		//

		//Increase Job COunter

		go func(j Job) {
			defer p.Wg.Done()
			p.JobMutex.Lock()
			p.JobCount += 1
			p.JobMutex.Unlock()
			//Execute the task
			er := j.Cx.Execute()

			//Add Error Handling If Defined
			if er != nil {
				if p.HandleError != nil {
					p.HandleError(er)
				}
			}
			//Tell Operator Task is completed and worker is free
			p.InterChannel <- JobResponse{
				Err: er,
				Uid: j.UID,
			}
			// ioutils.Cout.PrintInfo("completed %v\n", j.UID)
		}(rjob)

	}

	//When there is no work Shut it down
	close(p.InterChannel)
}

// Operator : (shared)Receives JobResponse and Frees the worker
func (p *Pool) Operator() {
	defer p.logic.Done()

	for {
		r, ok := <-p.InterChannel
		if !ok {
			break
		}

		//decrease active  job count
		p.JobMutex.Lock()
		p.JobCount -= 1

		// fmt.Printf("Job COunt %v\n", p.JobCount)
		if p.JobCount == 0 {
			// fmt.Println("sending all jobs completion")
			//send 1 to completion channel when all active jobs are completed
			p.Completion <- 1
		}

		p.JobMutex.Unlock()

		//If interchannel commuication has happened then work is completed
		//add worker back to queue
		p.workerChannel <- 8

		if p.OnCompletion != nil {
			p.OnCompletion(r)
		}

	}
}

// Done : Should Be called after all jobs are assigned
func (p *Pool) Done() {
	close(p.JobChannel)
}

// Release : Must be defered at beginning
func (p *Pool) Release() {
	defer p.logic.Wait()
}

// Wait : Wait For Ongoing Task to Complete
func (p *Pool) Wait() {
	// Do not wait if there are no jobs

	p.JobMutex.Lock()

	if p.JobCount == 0 && p.JobSent == 0 {
		p.JobMutex.Unlock()
		return
	}
	p.JobMutex.Unlock()
	<-p.Completion
	p.JobSent = 0
}

func NewPool(size int) *Pool {

	p := Pool{
		workerChannel: make(chan int, 4),
		JobChannel:    make(chan Job, 4),
		InterChannel:  make(chan JobResponse, 4),
		Completion:    make(chan int, 4),
		Concurrency:   size,
	}

	if p.Concurrency < 2 {
		p.Concurrency = 2
	}

	p.logic = &sync.WaitGroup{}
	p.Wg = &sync.WaitGroup{}
	p.JobMutex = &sync.Mutex{}

	p.logic.Add(2)
	go p.Launcher()
	go p.Operator()

	return &p
}
