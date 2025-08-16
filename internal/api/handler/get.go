package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"path"
	"strconv"

	"lo/domain/task"
	"lo/internal/api"
	"lo/internal/logger"
)

// GetTask godoc
// @Summary Get a task by ID
// @Description Получить задачу по ID.
// @Tags tasks
// @Produce application/json
// @Param id path int true "Task ID"
// @Success 200 {object} task.Task
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /tasks/{id} [get]
func GetTask(sc task.StorageClient, as *logger.AsyncLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, idStr := path.Split(r.URL.Path)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			api.WriteError(w, as, http.StatusBadRequest, "invalid id")
			as.Error("GetTask: invalid id", slog.String("id", idStr), slog.String("error", err.Error()))
			return
		}

		t := sc.Get(id)
		if t == nil {
			api.WriteError(w, as, http.StatusNotFound, "task not found")
			as.Error("GetTask: task not found", slog.Int("id", id))
			return
		}

		writeResponseWithTask(w, t, as)
		as.Info("GetTask: successfully get task", slog.Int("id", id))
	}
}

func writeResponseWithTask(w http.ResponseWriter, t *task.Task, as *logger.AsyncLogger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(t); err != nil {
		as.Error("writeResponseWIthTask: cannot encode request body", slog.String("error", err.Error()))
	}
}
