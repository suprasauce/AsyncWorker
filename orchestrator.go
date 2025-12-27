package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	Num int
}

type Orchestrator struct {
	jobQueue          chan Job
	wg                sync.WaitGroup
	mCtx              context.Context
	wCtx              context.Context
	wCtxCancel        context.CancelFunc
	jobQueueListener  chan Job
	currActiveWorkers int
	mWg               sync.WaitGroup
}

func (obj *Orchestrator) Init(jobQueueSize int, nWorkers int) {
	obj.mCtx = context.Background()
	obj.wCtx, obj.wCtxCancel = context.WithCancel(obj.mCtx)
	obj.jobQueue = make(chan Job, jobQueueSize)
	obj.jobQueueListener = make(chan Job) //  unbuffered
	// for i := 0; i < nWorkers; i++ {
	// 	go obj.SpoolWorker(obj.wCtx, i)
	// }
	go obj.JobQueueListenerService(obj.wCtx)
}

func (obj *Orchestrator) SpoolWorker(ctx context.Context, workerId int) {
	fmt.Printf("Worker with Id %d spooled\n", workerId)
	obj.wg.Add(1)
	defer func() {
		if ok := recover(); ok != nil {
			/* panic situaion arised, as of current implementation we will drop the current job
			   and respawn this worker
			*/
			fmt.Printf("Panic for worker Id: %d, message: %v\n", workerId, ok)
			obj.SpoolWorker(ctx, workerId)
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
	// for j := range obj.jobQueue {
	// 	obj.ProcessJob(ctx, workerId, j)
	// 	if j.Num == 211 {
	// 		pp := 0
	// 		a := 2 / pp // here panic will occur, so worker should respool
	// 		println(a)
	// 	}
	// }
	/*
		Above one doesnt use context, so using the below one
	*/
	for {
		select {
		case j, ok := <-obj.jobQueueListener:
			if !ok {
				return
			}
			obj.ProcessJob(ctx, workerId, j)
		case <-ctx.Done():
			return
		}
	}
}

func (obj *Orchestrator) JobQueueListenerService(ctx context.Context) {
	go func() {
		defer func() {
			obj.wCtxCancel()
			obj.wg.Wait()
			obj.mWg.Done()
		}()
		obj.mWg.Add(1)
		for j := range obj.jobQueue {
			select {
			case obj.jobQueueListener <- j:
			default:
				/*
					if we are here then this means that no worker is available, spool a new on if
					max limit is not reached.
				*/
				if obj.currActiveWorkers < N_WORKERS {
					// need to decide on the context though, whether it should be derived or not?
					go obj.SpoolWorker(ctx, rand.Int())
					obj.currActiveWorkers += 1
				}
				/*
					now the below code is controversial I think but it gets the work done.
					so lets we spool a worker then we need to send the job to the worker
					I can create a new function maybe like spoolworkerWithJob....
				*/
				obj.jobQueueListener <- j
			}
		}
	}()
}

func (obj *Orchestrator) SendJob(j Job) {
	obj.jobQueue <- j // this is currently blocking in nature
}

/*
we are passing worker context here, However we can create a job context derived from
worker context.
*/
func (obj *Orchestrator) ProcessJob(ctx context.Context, wokerId int, j Job) {
	time.Sleep(time.Second * 1)
	fmt.Printf("Work done by worker with Id %d, value is %d\n", wokerId, j.Num)
}

/*
For now cancellation policy is dropping jobs which are buffered but waiting for the current
processing jobs to be completed.
*/

/*
in our current implementation when we shutdown, no new jobs are taken but the jobs
present in the queue are still processed. After that the
*/
func (obj *Orchestrator) GraceFulShutdown() {
	close(obj.jobQueue)
	/*
		the below things are shifted to jobqueue listnere service.
	*/
	// obj.wCtxCancel() // cancelling all the workers
	// obj.wg.Wait()

	obj.mWg.Wait()
	fmt.Println("Async Worker Gracefully shutdown")
}

/*
we will create one job listener channel and one job queue channel
workers will listen to te job listner

*/
