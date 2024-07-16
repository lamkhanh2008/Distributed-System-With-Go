package worker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Job struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobChannel chan Job
type JobQueue chan chan Job
type Worker struct {
	Id      int
	JobChan JobChannel
	Queue   JobQueue
	Quit    chan struct{}
}

func New(id int, jobchan JobChannel, queue JobQueue, quit chan struct{}) *Worker {
	return &Worker{
		Id:      id,
		JobChan: jobchan,
		Queue:   queue,
		Quit:    quit,
	}
}

func (w *Worker) Start() {
	c := &http.Client{Timeout: time.Millisecond * 15000}
	go func() {
		for {
			w.Queue <- w.JobChan
			select {
			case job := <-w.JobChan:
				callApi(job.Id, w.Id, c)
			case <-w.Quit:
				close(w.JobChan)
				return
			}
		}
	}()

}

func (w *Worker) Stop() {
	close(w.Quit)
}

func callApi(num, id int, c *http.Client) {
	fmt.Printf("Job id %d served by Worker id %d", num, id)
	baseURL := "https://age-of-empires-2-api.herokuapp.com/api/v1/civilization/%d"

	ur := fmt.Sprintf(baseURL, num)
	req, err := http.NewRequest(http.MethodGet, ur, nil)
	if err != nil {
		//log.Printf("error creating a request for term %d :: error is %+v", num, err)
		return
	}
	res, err := c.Do(req)
	if err != nil {
		//log.Printf("error querying for term %d :: error is %+v", num, err)
		return
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		//log.Printf("error reading response body :: error is %+v", err)
		return
	}
	//log.Printf("%d  :: ok", id)
}
