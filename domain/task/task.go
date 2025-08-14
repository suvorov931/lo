package task

const (
	StatusNew     string = "new"
	StatusRunning string = "running"
	StatusDone    string = "done"
)

type Task struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
