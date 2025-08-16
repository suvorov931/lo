// @title Tasks API
// @version 1.0
// @description REST API для управления задачами (Task) с асинхронным логированием.
// @contact.name Your Name
// @contact.email you@example.com
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lo/domain/task"
	"lo/internal/api/handler"
	llogger "lo/internal/logger"
)

const (
	httpHost = "0.0.0.0"
	httpPort = 8080

	logChBuf = 100

	shutdownTime = 15 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	worker := llogger.NewWorker(logChBuf)
	logger := worker.AsyncLogger

	go worker.Run(ctx)

	storageTask := task.New()

	router := initRouter(storageTask, logger)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", httpHost, httpPort),
		Handler: router,
	}

	go func() {
		logger.Info("starting http server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", slog.String("err", err.Error()))
		}
	}()

	<-ctx.Done()

	logger.Info("received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTime)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("cannot shutdown http server", slog.String("err", err.Error()))
		return
	}

	logger.Info("stopping http server", slog.String("addr", server.Addr))

	logger.Info("application shutdown completed successfully")

	worker.Close(shutdownCtx)
}

func initRouter(storageTask *task.StorageTask, logger llogger.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateTask(storageTask, logger)(w, r)

		case http.MethodGet:
			handler.ListTasks(storageTask, logger)(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(storageTask, logger)(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	return mux
}
