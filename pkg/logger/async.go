package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

// asyncLogger wraps a slog.Logger and provides async buffered logging
type asyncLogger struct {
	logger  *slog.Logger
	queue   chan logEntry
	wg      sync.WaitGroup
	stop    chan struct{}
	started bool
	mu      sync.Mutex
}

type logEntry struct {
	level slog.Level
	msg   string
	args  []any
	ctx   context.Context
}

// defaultAsyncBufferSize is the default buffer size for async logging
const defaultAsyncBufferSize = 1000

var (
	asyncLog *asyncLogger
	asyncMu  sync.Mutex
)

// newAsyncLogger creates a new async logger wrapper
func newAsyncLogger(baseLogger *slog.Logger, bufferSize int) *asyncLogger {
	if bufferSize <= 0 {
		bufferSize = defaultAsyncBufferSize
	}

	al := &asyncLogger{
		logger:  baseLogger,
		queue:   make(chan logEntry, bufferSize),
		stop:    make(chan struct{}),
		started: true, // Set before starting goroutine to prevent race condition
	}

	al.wg.Add(1)
	go al.processLogs()

	return al
}

// processLogs processes log entries from the queue
func (al *asyncLogger) processLogs() {
	defer al.wg.Done()

	for {
		select {
		case entry := <-al.queue:
			// Note: slog.Logger.Log() doesn't return an error.
			// Error handling should be done at the handler level if needed.
			al.logger.Log(entry.ctx, entry.level, entry.msg, entry.args...)
		case <-al.stop:
			// Drain remaining entries
			for {
				select {
				case entry := <-al.queue:
					// Note: slog.Logger.Log() doesn't return an error.
					// Error handling should be done at the handler level if needed.
					al.logger.Log(entry.ctx, entry.level, entry.msg, entry.args...)
				default:
					return
				}
			}
		}
	}
}

// log enqueues a log entry
func (al *asyncLogger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	entry := logEntry{
		level: level,
		msg:   msg,
		args:  args,
		ctx:   ctx,
	}

	select {
	case al.queue <- entry:
		// Successfully queued
	default:
		// Queue is full, log synchronously to prevent blocking
		// Note: slog.Logger.Log() doesn't return an error.
		// Error handling should be done at the handler level if needed.
		al.logger.Log(ctx, level, msg, args...)
	}
}

// Debug logs a debug message asynchronously
func (al *asyncLogger) Debug(ctx context.Context, msg string, args ...any) {
	if al.logger.Enabled(ctx, slog.LevelDebug) {
		al.log(ctx, slog.LevelDebug, msg, args...)
	}
}

// Info logs an info message asynchronously
func (al *asyncLogger) Info(ctx context.Context, msg string, args ...any) {
	al.log(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs a warning message asynchronously
func (al *asyncLogger) Warn(ctx context.Context, msg string, args ...any) {
	al.log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs an error message asynchronously
func (al *asyncLogger) Error(ctx context.Context, msg string, args ...any) {
	al.log(ctx, slog.LevelError, msg, args...)
}

// Shutdown gracefully shuts down the async logger
func (al *asyncLogger) Shutdown() {
	al.mu.Lock()
	if !al.started {
		al.mu.Unlock()
		return
	}
	al.started = false
	al.mu.Unlock()

	// Signal shutdown
	close(al.stop)

	// Wait for goroutine to finish with timeout
	done := make(chan struct{})
	go func() {
		al.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Successfully drained
	case <-time.After(5 * time.Second):
		// Timeout - log warning but continue
		fmt.Fprintf(os.Stderr, "Async logger shutdown timeout: some logs may be lost\n")
	}
}

// enableAsync enables async logging for the global logger
func enableAsync(baseLogger *slog.Logger, bufferSize int) {
	asyncMu.Lock()
	defer asyncMu.Unlock()

	if asyncLog != nil {
		asyncLog.Shutdown()
	}

	asyncLog = newAsyncLogger(baseLogger, bufferSize)
	// started is now set in newAsyncLogger before goroutine starts
}

// disableAsync disables async logging
func disableAsync() {
	asyncMu.Lock()
	defer asyncMu.Unlock()

	if asyncLog != nil {
		asyncLog.Shutdown()
		asyncLog = nil
	}
}
