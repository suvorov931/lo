package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"lo/domain/task"
	"lo/internal/logger"
)

func CreateTask(st *task.StorageTask, as *logger.AsyncLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			as.Error(ctx, "CreateTask: cannot decode request body", slog.String("error", err.Error()))

			return
		}

		st.Save(&t)

		writeResponse(ctx, w, t.Id, as)

		as.Info(ctx, "CreateTask: successfully created task", slog.Int("id", t.Id))
	}
}

func writeResponse(ctx context.Context, w http.ResponseWriter, id int, as *logger.AsyncLogger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := struct {
		Id int `json:"id"`
	}{
		Id: id,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		as.Error(ctx, "CreateTask: cannot decode request body", slog.String("error", err.Error()))
	}
}
