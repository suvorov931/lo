package task

import (
	"sync"
)

type Task struct {
	Id      int
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type StorageTask struct {
	countId int
	mu      *sync.Mutex
	storage map[int]*Task
}

func New() *StorageTask {
	return &StorageTask{
		countId: 0,
		mu:      &sync.Mutex{},
		storage: make(map[int]*Task),
	}
}

func (st *StorageTask) Save(task *Task) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.countId++
	task.Id = st.countId

	st.storage[st.countId] = task
}
