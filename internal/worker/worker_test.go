package worker

import (
	"sync"
	"testing"

	"github.com/rcbadiale/go-stress-test/internal/tasks"
	"github.com/stretchr/testify/suite"
)

type WorkerTestSuite struct {
	taskMock *tasks.TaskMock
	toDo     chan tasks.Task
	done     chan tasks.Result
	wg       *sync.WaitGroup
	suite.Suite
}

func (s *WorkerTestSuite) SetupTest() {
	s.taskMock = new(tasks.TaskMock)
	s.toDo = make(chan tasks.Task, 1)
	s.done = make(chan tasks.Result, 1)
	s.wg = &sync.WaitGroup{}
}

func (s *WorkerTestSuite) TearDownTest() {
	close(s.toDo)
	close(s.done)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(WorkerTestSuite))
}

func (s *WorkerTestSuite) TestRunGivenTaskOnToDoChannelShouldWriteResultToDoneChannel() {
	s.taskMock.On("Execute").Return(
		tasks.Result{TaskName: "TaskMock", Status: "200"},
	)

	worker := NewWorker(1, s.wg, s.toDo, s.done)
	s.wg.Add(1)
	go worker.Run()
	s.toDo <- s.taskMock
	s.wg.Wait()
	s.taskMock.AssertExpectations(s.T())

	result := <-s.done
	s.Equal("TaskMock", result.TaskName)
	s.Equal("200", result.Status)
}
