package task

import "sync"

type RequestTask struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Task struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type StorageTask struct {
	countId int
	mu      *sync.Mutex
	storage map[int]*Task
}

type StorageClient interface {
	Save(task *Task)
	Get(id int) *Task
	GetAll() []*Task
	GetByStatus(status string) []*Task
}
