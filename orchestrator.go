package main

import (
	"fmt"
	"sync"
)

type Job struct {
	Num int
}

type Orchestrator struct {
	jobQueue chan Job
	wg       sync.WaitGroup
}

func (obj *Orchestrator) Init(jobQueueSize int, nWorkers int) {
	obj.jobQueue = make(chan Job, jobQueueSize)
	for i := 0;i<nWorkers;i++ {
		go obj.SpoolWorker(i)
	}
}

func (obj *Orchestrator) SpoolWorker(workerId int) {
	lg.SendLog(fmt.Sprintf("Worker with Id %d spooled", workerId))
	obj.wg.Add(1)
	defer func() {
		if ok := recover(); ok != nil {
			/* panic situaion arised, as of current implementation we will drop the current job
			   and respawn this worker
			*/
			lg.SendLog(fmt.Sprintf("Panic for worker Id: %d, message: %v", workerId, ok))
			obj.SpoolWorker(workerId)
		}
		obj.wg.Done()
	}()
	// for {
	// 	j, ok := <- obj.jobQueue
	// 	if !ok {
	// 		return
	// 	}
	// 	obj.ProcessJob(j)
	// }
	/*
	Below is more idiomatic. The range automatically handles if channel is closed or not
	and exits the loop when there are no more values to read.
	*/
	for j := range obj.jobQueue {
		obj.ProcessJob(workerId, j)
		if j.Num == 211 {
			pp := 0
			a := 2/pp // here panic will occur, so worker should respool
			println(a)
		}
	}
}

func (obj *Orchestrator) SendJob(j Job) {
	obj.jobQueue <- j
}

func (obj *Orchestrator) ProcessJob(wokerId int, j Job) {
	lg.SendLog(fmt.Sprintf("Work done by worker with Id %d, value is %d", wokerId, j.Num))
}

func (obj *Orchestrator) GraceFulShutdown() {
	close(obj.jobQueue)
	obj.wg.Wait()
	lg.SendLog(fmt.Sprintf("Async Worker Gracefully shutdown"))
}


