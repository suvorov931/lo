package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"lo/domain/task"
	"lo/internal/api"
	"lo/internal/logger"
)

func CreateTask(st *task.StorageTask, as *logger.AsyncLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			api.WriteError(w, as, http.StatusBadRequest, "invalid body")
			as.Error("CreateTask: cannot decode request body", slog.String("error", err.Error()))
			return
		}

		st.Save(&t)

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
