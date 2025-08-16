package logger

import (
	"log/slog"
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

type Worker struct {
	AsyncLogger *AsyncLogger
	wg          *sync.WaitGroup
}

type Logger interface {
	Info(message string, attrs ...slog.Attr)
	Error(message string, attrs ...slog.Attr)
	Warn(message string, attrs ...slog.Attr)
}
