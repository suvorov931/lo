package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"lo/domain/task"
	"lo/internal/logger"
)

func CreateTask(logCh chan logger.Log) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			logCh <- logger.Log{
				Level:   slog.LevelError,
				Message: "CreateTask: cannot decode request body",
			}
		}

		logCh <- logger.Log{
			Level:   slog.LevelInfo,
			Message: "hello where",
			Attrs:   []slog.Attr{{Key: "1", Value: slog.Value{}}},
		}
		logCh <- logger.Log{
			Level:   slog.LevelInfo,
			Message: "salam where",
		}
		fmt.Println(t)
	}
}
