package manager

import (
	"testing"

	"github.com/rcbadiale/go-stress-test/internal/tasks"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	taskMock *tasks.TaskMock
	suite.Suite
}

func (s *ManagerTestSuite) SetupTest() {
	s.taskMock = new(tasks.TaskMock)
}

func (s *ManagerTestSuite) TearDownTest() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}

func (s *ManagerTestSuite) TestPublishGivenTasksShouldWriteToToDoChannel() {
	tasks := []tasks.Task{}
	for range 5 {
		tasks = append(tasks, s.taskMock)
	}
	manager := NewManager(1, tasks)
	manager.publish()

	s.Len(manager.toDo, 5)
	s.Equal(s.taskMock, <-manager.toDo)
}

func (s *ManagerTestSuite) TestExecuteGivenTasksShouldReturnDuration() {
	s.taskMock.On("Execute").Return(
		tasks.Result{TaskName: "TaskMock", Status: "200"},
	)

	tasks := []tasks.Task{s.taskMock}
	manager := NewManager(1, tasks)
	duration := manager.Execute()

	s.taskMock.AssertExpectations(s.T())
	s.NotZero(duration)
}

func (s *ManagerTestSuite) TestReportGivenTasksShouldReturnReport() {
	s.taskMock.On("Execute").Return(
		tasks.Result{TaskName: "TaskMock", Status: "DONE"},
	)

	tasks := []tasks.Task{s.taskMock}
	manager := NewManager(1, tasks)
	manager.Execute()
	report := manager.Report()

	s.taskMock.AssertExpectations(s.T())
	s.Contains(report, "Results:\n- DONE: 1/1 (100%)\nTotal tasks executed: 1")
}
