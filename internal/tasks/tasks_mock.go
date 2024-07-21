package tasks

import "github.com/stretchr/testify/mock"

type TaskMock struct {
	mock.Mock
}

func (m *TaskMock) Name() string {
	return "TaskMock"
}

func (m *TaskMock) Execute() Result {
	args := m.Called()
	return args.Get(0).(Result)
}
