package logger

import (
	"context"
	"log/slog"
	"os"
)

const logChBuf = 100

func NewLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)

	return logger
}

type Log struct {
	Level   slog.Level
	Message string
	Attrs   []slog.Attr
}

type AsyncLogger struct {
	Logger *slog.Logger
	LogCh  chan Log
}

func NewLoggingWorker(logger *slog.Logger) *AsyncLogger {
	return &AsyncLogger{
		Logger: logger,
		LogCh:  make(chan Log, logChBuf),
	}
}

func (lw *AsyncLogger) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case log, ok := <-lw.LogCh:
			if !ok {
				return
			}

			lw.Logger.LogAttrs(ctx, log.Level, log.Message, log.Attrs...)
		}
	}
}

func (lw *AsyncLogger) Close() {
	close(lw.LogCh)
}
