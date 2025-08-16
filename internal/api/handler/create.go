package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"lo/domain/task"
	"lo/internal/api"
	"lo/internal/logger"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Создать задачу. ID генерируется сервером.
// @Tags tasks
// @Accept application/json
// @Produce application/json
// @Param task body task.RequestTask true "Task to create"      // model package path: domain/task
// @Success 201 {object} map[string]int "created id"     // можно описать структуру ответа
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /tasks [post]
func CreateTask(sc task.StorageClient, as *logger.AsyncLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			api.WriteError(w, as, http.StatusBadRequest, "invalid body")
			as.Error("CreateTask: cannot decode request body", slog.String("error", err.Error()))
			return
		}

		sc.Save(&t)

		writeResponse(w, t.Id, as)

		as.Info("CreateTask: successfully created task", slog.Int("id", t.Id))
	}
}

func writeResponse(w http.ResponseWriter, id int, as *logger.AsyncLogger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := struct {
		Id int `json:"id"`
	}{
		Id: id,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		as.Error("CreateTask: cannot decode request body", slog.String("error", err.Error()))
	}
}
