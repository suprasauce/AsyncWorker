package main

type Job struct {
	Num int
}

var jobQueue chan Job
const JOB_QUEUE_SIZE int = 4
const N_WORKERS int = 2


func SpoolWorker() {
	for {
		j:=  <- jobQueue
		println(j.Num)
	}
}

func init() {
	jobQueue = make(chan Job, JOB_QUEUE_SIZE)
	for i := 0;i<N_WORKERS;i++ {
		go SpoolWorker()
	}
}

func SendJob(j Job) {
	jobQueue <- j // currently this will block when queue is full
}

func main() {
	for i := 0;i<10;i++ {
		SendJob(Job{Num: i})
	}
	for {

	}
}