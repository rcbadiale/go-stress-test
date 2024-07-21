package worker

import (
	"sync"

	"github.com/rcbadiale/go-stress-test/internal/tasks"
)

type Worker struct {
	Id      int
	wg      *sync.WaitGroup
	toDo    <-chan tasks.Task
	results chan<- tasks.Result
}

func NewWorker(id int, wg *sync.WaitGroup, toDo <-chan tasks.Task, results chan<- tasks.Result) *Worker {
	return &Worker{Id: id, wg: wg, toDo: toDo, results: results}
}

func (w *Worker) Run() {
	for job := range w.toDo {
		result := job.Execute()
		result.WorkerId = w.Id
		w.results <- result
		w.wg.Done()
	}
}
