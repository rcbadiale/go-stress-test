package tasks

import (
	"fmt"
	"io"
	"net/http"
)

type RequestTask struct {
	id     int
	client *http.Client
	req    *http.Request
	resp   *http.Response
}

func NewRequestTask(id int, client *http.Client, method, url string, body io.Reader) (*RequestTask, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("cannot create task: %s", err)
	}
	return &RequestTask{id: id, client: client, req: req}, nil
}

func (r *RequestTask) Name() string {
	return fmt.Sprintf("RequestTask-%d", r.id)
}

func (r *RequestTask) Execute() Result {
	res := Result{TaskName: r.Name()}
	resp, err := r.client.Do(r.req)
	if err != nil {
		res.Status = err.Error()
		return res
	}
	r.resp = resp
	res.Status = fmt.Sprintf("%d", resp.StatusCode)
	return res
}
