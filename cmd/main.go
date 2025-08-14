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

	"lo/internal/api/handler"
	llogger "lo/internal/logger"
)

const (
	httpHost = "localhost"
	httpPort = 8080

	shutdownTime = 15 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	logger := llogger.NewLogger()

	loggingWorker := llogger.NewLoggingWorker(logger)

	go loggingWorker.Run(ctx)

	router := initRouter(loggingWorker)

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

	loggingWorker.Close()

	logger.Info("stopping http server", slog.String("addr", server.Addr))

	logger.Info("application shutdown completed successfully")
}

func initRouter(loggingWorker *llogger.AsyncLogger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateTask(loggingWorker.LogCh)(w, r)

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
