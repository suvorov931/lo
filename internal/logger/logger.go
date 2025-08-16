package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

func NewWorker(logChBuf int) *Worker {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)

	logCh := make(chan logEntry, logChBuf)

	wg := &sync.WaitGroup{}

	return &Worker{
		AsyncLogger: &AsyncLogger{
			logger: logger,
			logCh:  logCh,
			wg:     wg,
		},
		wg: wg,
	}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for log := range w.AsyncLogger.logCh {
				w.AsyncLogger.logger.LogAttrs(context.Background(), log.level, log.message, log.attrs...)
				w.wg.Done()
			}

			return

		case log, ok := <-w.AsyncLogger.logCh:
			if !ok {
				return
			}

			w.AsyncLogger.logger.LogAttrs(ctx, log.level, log.message, log.attrs...)
			w.wg.Done()
		}
	}
}

func (al *AsyncLogger) Info(message string, attrs ...slog.Attr) {
	al.wg.Add(wgDefaultDelta)

	select {
	case al.logCh <- logEntry{
		level:   slog.LevelInfo,
		message: message,
		attrs:   attrs,
	}:
	default:
		al.wg.Done()
	}

}

func (al *AsyncLogger) Error(message string, attrs ...slog.Attr) {
	al.wg.Add(wgDefaultDelta)

	select {
	case al.logCh <- logEntry{
		level:   slog.LevelError,
		message: message,
		attrs:   attrs,
	}:
	default:
		al.wg.Done()
	}
}

func (al *AsyncLogger) Warn(message string, attrs ...slog.Attr) {
	al.wg.Add(wgDefaultDelta)

	select {
	case al.logCh <- logEntry{
		level:   slog.LevelWarn,
		message: message,
		attrs:   attrs,
	}:
	default:
		al.wg.Done()
	}
}

func (w *Worker) Close(ctx context.Context) {
	done := make(chan struct{})

	go func() {
		close(w.AsyncLogger.logCh)
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		w.AsyncLogger.logger.LogAttrs(context.Background(), slog.LevelWarn,
			"AsyncLogger: close timeout reached, some logs may be lost",
			slog.String("reason", ctx.Err().Error()),
		)
	}
}
