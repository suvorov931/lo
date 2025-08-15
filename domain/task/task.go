package task

import (
	"fmt"
	"sync"
)

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

func (st *StorageTask) Get(id int) (*Task, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	task, ok := st.storage[id]
	if !ok {
		return nil, fmt.Errorf("Get: task %d not found", id)
	}

	return task, nil
}
