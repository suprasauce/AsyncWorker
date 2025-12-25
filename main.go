package main

const JOB_QUEUE_SIZE int = 4
const N_WORKERS int = 2

var orch Orchestrator
var lg logger

func init() {
	orch.Init(JOB_QUEUE_SIZE, N_WORKERS)
	lg.Init()
}

func main() {
	for i := 0; i < 10; i++ {
		orch.SendJob(Job{Num: i})
	}
	orch.GraceFulShutdown()
	lg.GraceFulShutdown()
}