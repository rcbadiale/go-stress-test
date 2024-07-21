package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/rcbadiale/go-stress-test/internal/tasks"
	"github.com/rcbadiale/go-stress-test/internal/worker"
)

type Manager struct {
	workers []*worker.Worker
	tasks   []tasks.Task
	wg      *sync.WaitGroup
	toDo    chan tasks.Task
	done    chan tasks.Result
}

func NewManager(workersCount int, t []tasks.Task) *Manager {
	tasksLen := len(t)
	toDo := make(chan tasks.Task, tasksLen)
	done := make(chan tasks.Result, tasksLen)
	wg := &sync.WaitGroup{}
	return &Manager{
		workers: spawnWorkers(workersCount, wg, toDo, done),
		tasks:   t,
		toDo:    toDo,
		done:    done,
		wg:      wg,
	}
}

func (m *Manager) publish() {
	defer close(m.toDo)
	for _, task := range m.tasks {
		m.wg.Add(1)
		m.toDo <- task
	}
}

func (m *Manager) Execute() time.Duration {
	for _, worker := range m.workers {
		go worker.Run()
	}
	start := time.Now()
	m.publish()
	m.wg.Wait()
	close(m.done)
	return time.Since(start)
}

func (m *Manager) Report() string {
	results := make(map[string]int)
	total := 0

	for result := range m.done {
		total += 1
		if result.Status != "" {
			if _, ok := results[result.Status]; !ok {
				results[result.Status] = 0
			}
			results[result.Status] += 1
		}
	}
	output := "Results:\n"
	for key, val := range results {
		output += fmt.Sprintf("- %s: %d/%d (%d%%)\n", key, val, total, 100*val/total)
	}
	output += fmt.Sprintf("Total tasks executed: %d", total)
	return output
}

func spawnWorkers(workersCount int, wg *sync.WaitGroup, toDo <-chan tasks.Task, done chan<- tasks.Result) []*worker.Worker {
	workers := []*worker.Worker{}
	for i := range workersCount {
		workers = append(workers, worker.NewWorker(i, wg, toDo, done))
	}
	return workers
}
