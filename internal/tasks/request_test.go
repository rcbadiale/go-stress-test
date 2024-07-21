package tasks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RequestTaskTestSuite struct {
	mockServer *httptest.Server
	mockClient *http.Client
	badServer  *httptest.Server
	badClient  *http.Client
	suite.Suite
}

func (s *RequestTaskTestSuite) SetupTest() {
	s.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	s.mockClient = s.mockServer.Client()
	s.badServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	}))
	s.badClient = s.badServer.Client()
}

func (s *RequestTaskTestSuite) TearDownTest() {
	s.mockServer.Close()
	s.badServer.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RequestTaskTestSuite))
}

func (s *RequestTaskTestSuite) TestNameShouldReturnTaskName() {
	task := RequestTask{id: 999}
	s.Equal("RequestTask-999", task.Name())
}

func (s *RequestTaskTestSuite) TestNewRequestGivenBadDataShouldFail() {
	_, err := NewRequestTask(1, s.mockClient, "GET", "://", nil)
	s.ErrorContains(err, "cannot create task")
}

func (s *RequestTaskTestSuite) TestExecuteGivenSuccessShouldReturnResultWithStatusCode() {
	task, err := NewRequestTask(1, s.mockClient, "GET", s.mockServer.URL, nil)
	s.NoError(err)
	result := task.Execute()
	s.Equal("200", result.Status)
	s.Equal("RequestTask-1", result.TaskName)
}

func (s *RequestTaskTestSuite) TestExecuteGivenFailureShouldReturnResultWithStatusCode() {
	task, err := NewRequestTask(2, s.badClient, "GET", s.badServer.URL, nil)
	s.NoError(err)
	result := task.Execute()
	errorMessage := fmt.Sprintf("Get \"%s\": EOF", s.badServer.URL)
	s.Equal(errorMessage, result.Status)
	s.Equal("RequestTask-2", result.TaskName)
}
