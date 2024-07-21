package tasks

type Task interface {
	Name() string
	Execute() Result
}

type Result struct {
	WorkerId int
	TaskName string
	Status   string
}
