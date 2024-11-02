package dispatcher

import "worker_pool/worker"

type dispatcher struct {
	workers []*worker.Worker
	JobChan worker.JobChannel
	Queue   worker.JobQueue
}

func New() *dispatcher {

	return &dispatcher{
		workers: make([]*worker.Worker, 16),
		Queue:   make(worker.JobQueue),
	}

}

func (dis *dispatcher) Start() *dispatcher {
	l := len(dis.workers)
	for i := 0; i < l; i++ {
		w := worker.New(i, make(worker.JobChannel), dis.Queue, make(chan struct{}))
		w.Start()
	}
	go dis.process()
	return dis
}

func (dis *dispatcher) process() {
	for {
		select {
		case job := <-dis.JobChan:
			jobChan := <-dis.Queue
			jobChan <- job
		}
	}
}

func (dis *dispatcher) SubmitJob(job worker.Job) {
	dis.JobChan <- job
}
