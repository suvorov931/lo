package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"lo/domain/task"
	"lo/internal/api"
	"lo/internal/logger"
)

// ListTasks godoc
// @Summary List tasks
// @Description Получить список задач. Можно фильтровать по статусу через query param ?status=<status>.
// @Tags tasks
// @Produce application/json
// @Param status query string false "Filter tasks by status (optional)"
// @Success 200 {array} task.Task
// @Failure 500 {object} api.ErrorResponse
// @Router /tasks [get]
func ListTasks(sc task.StorageClient, logger logger.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []*task.Task

		status := r.URL.Query().Get("status")
		if status == "" {
			tasks = sc.GetAll()
		} else {
			tasks = sc.GetByStatus(status)
		}

		if tasks == nil {
			api.WriteError(w, logger, http.StatusNotFound, "tasks not found")
			logger.Warn("ListTasks: tasks not found")
			return
		}

		writeResponseWithTasks(w, tasks, logger)
		logger.Info("ListTasks: successfully list tasks", slog.Int("count", len(tasks)))
	}
}

func writeResponseWithTasks(w http.ResponseWriter, t []*task.Task, logger logger.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(t); err != nil {
		logger.Error("writeResponseWithTask: cannot encode request body", slog.String("error", err.Error()))
	}
}
