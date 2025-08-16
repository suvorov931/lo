package task

import (
	"sync"
)

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

func (st *StorageTask) Get(id int) *Task {
	st.mu.Lock()
	defer st.mu.Unlock()

	task, ok := st.storage[id]
	if !ok {
		return nil
	}

	return task
}

func (st *StorageTask) GetAll() []*Task {
	st.mu.Lock()
	defer st.mu.Unlock()

	var tasks []*Task

	for _, task := range st.storage {
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil
	}

	return tasks
}

func (st *StorageTask) GetByStatus(status string) []*Task {
	st.mu.Lock()
	defer st.mu.Unlock()

	var tasks []*Task

	for _, task := range st.storage {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
