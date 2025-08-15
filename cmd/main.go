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
	httpHost = "localhost"
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

	asyncLogger := llogger.NewAsyncLogger(logChBuf)

	go asyncLogger.Run(ctx)

	storageTask := task.New()

	router := initRouter(storageTask, asyncLogger)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", httpHost, httpPort),
		Handler: router,
	}

	go func() {
		asyncLogger.Info(ctx, "starting http server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			asyncLogger.Error(ctx, "failed to start server", slog.String("err", err.Error()))
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTime)
	defer shutdownCancel()

	asyncLogger.Info(shutdownCtx, "received shutdown signal")

	if err := server.Shutdown(shutdownCtx); err != nil {
		asyncLogger.Error(shutdownCtx, "cannot shutdown http server", slog.String("err", err.Error()))
		return
	}

	asyncLogger.Info(shutdownCtx, "stopping http server", slog.String("addr", server.Addr))

	asyncLogger.Info(shutdownCtx, "application shutdown completed successfully")

	asyncLogger.Close(shutdownCtx)
}

func initRouter(storageTask *task.StorageTask, loggingWorker *llogger.AsyncLogger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateTask(storageTask, loggingWorker)(w, r)

		case http.MethodGet:
			handler.GetTask()(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListTasks()(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	return mux
}
