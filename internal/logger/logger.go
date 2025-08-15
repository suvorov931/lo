package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

const wgDefaultDelta = 1

type logEntry struct {
	level   slog.Level
	message string
	attrs   []slog.Attr
}

type AsyncLogger struct {
	logger *slog.Logger
	logCh  chan logEntry
	wg     *sync.WaitGroup
}

func NewAsyncLogger(logChBuf int) *AsyncLogger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)

	return &AsyncLogger{
		logger: logger,
		logCh:  make(chan logEntry, logChBuf),
		wg:     &sync.WaitGroup{},
	}
}

func (al *AsyncLogger) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for log := range al.logCh {
				al.logger.LogAttrs(context.Background(), log.level, log.message, log.attrs...)
				al.wg.Done()
			}

			return

		case log, ok := <-al.logCh:
			if !ok {
				return
			}

			al.logger.LogAttrs(ctx, log.level, log.message, log.attrs...)
			al.wg.Done()
		}
	}
}

func (al *AsyncLogger) Info(message string, attrs ...slog.Attr) {
	select {
	case al.logCh <- logEntry{
		level:   slog.LevelInfo,
		message: message,
		attrs:   attrs,
	}:
		al.wg.Add(wgDefaultDelta)

	default:

	}
}

func (al *AsyncLogger) Error(message string, attrs ...slog.Attr) {
	select {

	case al.logCh <- logEntry{
		level:   slog.LevelError,
		message: message,
		attrs:   attrs,
	}:
		al.wg.Add(wgDefaultDelta)

	default:

	}
}

func (al *AsyncLogger) Warn(message string, attrs ...slog.Attr) {
	select {
	case al.logCh <- logEntry{
		level:   slog.LevelWarn,
		message: message,
		attrs:   attrs,
	}:
		al.wg.Add(wgDefaultDelta)

	default:

	}
}

func (al *AsyncLogger) Close(ctx context.Context) {
	done := make(chan struct{})

	go func() {
		close(al.logCh)
		al.wg.Wait()

		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		al.logger.LogAttrs(context.Background(), slog.LevelWarn,
			"AsyncLogger: close timeout reached, some logs may be lost",
			slog.String("reason", ctx.Err().Error()),
		)
	}
}
