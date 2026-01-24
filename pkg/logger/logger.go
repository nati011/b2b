package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"marketplace/internal/config"
)

// Logger is the global logger instance
var Logger *slog.Logger

// appContext holds application-level context for logging
var (
	appContext map[string]any
	appCtxMu   sync.RWMutex
)

// useAsync indicates if async logging is enabled
var useAsync atomic.Bool

// sanitizeEnabled indicates if sanitization is enabled
var sanitizeEnabled atomic.Bool

// IsProduction returns true if the current environment is production
func IsProduction() bool {
	appCtxMu.RLock()
	defer appCtxMu.RUnlock()

	if appContext == nil {
		return false
	}

	// Check ECS field name first, then fallback to old name
	env, ok := appContext[ecsServiceEnvironment].(string)
	if !ok {
		env, ok = appContext["env"].(string)
	}
	return ok && strings.ToLower(env) == "production"
}

// Init initializes the global logger based on configuration
func Init(level string, format string) {
	InitWithConfig(level, format, false, 1.0, false, false, nil)
}

// Default initializes the logger with default settings (info level, text format)
func Default() {
	Init("info", "text")
}

// InitWithConfig initializes the logger with full configuration
func InitWithConfig(level string, format string, sourceLocation bool, samplingRate float64, async bool, sanitize bool, appCfg *config.AppConfig) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: sourceLocation,
	}

	if format == "json" {
		handler = newECSHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)

	// Set application context (using ECS field names)
	appCtxMu.Lock()
	if appCfg != nil {
		appContext = map[string]any{
			ecsServiceName:        appCfg.Name,
			ecsServiceVersion:     appCfg.Version,
			ecsServiceEnvironment: appCfg.Env,
			// Keep old names for backward compatibility during transition
			"app_name":    appCfg.Name,
			"app_version": appCfg.Version,
			"env":         appCfg.Env,
		}
	} else {
		appContext = make(map[string]any)
	}
	appCtxMu.Unlock()

	// Set sampling rate
	setSamplingRate(samplingRate)

	// Enable async logging if requested
	useAsync.Store(async)
	if async {
		enableAsync(Logger, defaultAsyncBufferSize)
	} else {
		disableAsync()
	}

	// Enable sanitization if requested
	sanitizeEnabled.Store(sanitize)
}

// WithAppContext sets application context for all logs
func WithAppContext(cfg *config.Config) {
	if cfg != nil {
		appCtxMu.Lock()
		appContext = map[string]any{
			ecsServiceName:        cfg.App.Name,
			ecsServiceVersion:     cfg.App.Version,
			ecsServiceEnvironment: cfg.App.Env,
			// Keep old names for backward compatibility during transition
			"app_name":    cfg.App.Name,
			"app_version": cfg.App.Version,
			"env":         cfg.App.Env,
		}
		appCtxMu.Unlock()
	}
}

// enrichArgs adds application context and sanitizes arguments
func enrichArgs(args ...any) []any {
	appCtxMu.RLock()
	ctxCopy := make(map[string]any, len(appContext))
	for k, v := range appContext {
		ctxCopy[k] = v
	}
	appCtxMu.RUnlock()

	enriched := make([]any, 0, len(args)+len(ctxCopy)*2+10)

	// Add ECS standard fields first
	enriched = addECSFields(enriched, ctxCopy)

	// Add app context (for backward compatibility)
	for k, v := range ctxCopy {
		// Skip if already added as ECS field
		if strings.HasPrefix(k, "service.") || k == "host.name" || k == "ecs.version" {
			continue
		}
		enriched = append(enriched, k, v)
	}

	// Add sanitized user args
	enriched = append(enriched, Sanitize(args...)...)
	return enriched
}

// logInternal is the internal logging function that handles async and sampling
func logInternal(ctx context.Context, level slog.Level, msg string, args ...any) {
	if Logger == nil {
		return
	}

	// Check if logging is enabled for this level
	if !Logger.Enabled(ctx, level) {
		return
	}

	// Apply sampling (errors always logged)
	rate := getSamplingRate()
	if !shouldSample(level, rate) {
		return
	}

	enrichedArgs := enrichArgs(args...)

	// Use async logger if enabled
	if useAsync.Load() && asyncLog != nil {
		asyncLog.log(ctx, level, msg, enrichedArgs...)
		return
	}

	// Synchronous logging
	// Note: slog.Logger.Log() doesn't return an error.
	// Error handling should be done at the handler level if needed.
	Logger.Log(ctx, level, msg, enrichedArgs...)
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	logInternal(context.Background(), slog.LevelDebug, msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	logInternal(context.Background(), slog.LevelInfo, msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	logInternal(context.Background(), slog.LevelWarn, msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	logInternal(context.Background(), slog.LevelError, msg, args...)
}

// Fatal logs a fatal message and exits
// Note: Fatal always logs synchronously to ensure the message is written before exit
func Fatal(msg string, args ...any) {
	if Logger != nil {
		enrichedArgs := enrichArgs(args...)
		Logger.Error(msg, enrichedArgs...)
	}
	os.Exit(1)
}

// Panic logs a panic message and panics
// Note: Panic always logs synchronously to ensure the message is written before panic
func Panic(msg string, args ...any) {
	if Logger != nil {
		enrichedArgs := enrichArgs(args...)
		Logger.Error(msg, enrichedArgs...)
	}
	panic(fmt.Sprintf(msg, args...))
}

// Context-aware logging functions

// WithContext creates a logger with context attributes
func WithContext(ctx context.Context) *slog.Logger {
	if Logger == nil {
		return nil
	}

	attrs := extractContextAttrs(ctx)
	if len(attrs) == 0 {
		return Logger
	}

	// Convert slog.Attr to slog.LogValuer-compatible format
	logger := Logger
	for _, attr := range attrs {
		logger = logger.With(attr.Key, attr.Value.Any())
	}
	return logger
}

// FromContext creates a logger from context with all context values
// Alias for WithContext for semantic clarity.
// Both FromContext and WithContext are available - use whichever is more semantically appropriate for your use case.
func FromContext(ctx context.Context) *slog.Logger {
	return WithContext(ctx)
}

// Context-aware logging functions that use the internal logging with async and sampling

// DebugContext logs a debug message with context
func DebugContext(ctx context.Context, msg string, args ...any) {
	logInternal(ctx, slog.LevelDebug, msg, args...)
}

// InfoContext logs an info message with context
func InfoContext(ctx context.Context, msg string, args ...any) {
	logInternal(ctx, slog.LevelInfo, msg, args...)
}

// WarnContext logs a warning message with context
func WarnContext(ctx context.Context, msg string, args ...any) {
	logInternal(ctx, slog.LevelWarn, msg, args...)
}

// ErrorContext logs an error message with context
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logInternal(ctx, slog.LevelError, msg, args...)
}

// extractContextAttrs extracts logging attributes from context
func extractContextAttrs(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr

	// Request/Trace ID (ECS: trace.id)
	if reqID := RequestIDFromContext(ctx); reqID != "" {
		attrs = append(attrs, slog.String("trace.id", reqID))
		// Keep old name for backward compatibility
		attrs = append(attrs, slog.String("request_id", reqID))
	}

	// Trace ID (ECS: trace.id)
	if traceID := TraceIDFromContext(ctx); traceID != "" {
		attrs = append(attrs, slog.String("trace.id", traceID))
		// Keep old name for backward compatibility
		attrs = append(attrs, slog.String("trace_id", traceID))
	}

	// User ID (ECS: user.id)
	if userID := UserIDFromContext(ctx); userID != "" {
		attrs = append(attrs, slog.String("user.id", userID))
		// Keep old name for backward compatibility
		attrs = append(attrs, slog.String("user_id", userID))
	}

	return attrs
}

// Helper functions for structured logging

// WithError adds an error to log arguments (ECS: error.message)
func WithError(err error) []any {
	if err == nil {
		return []any{}
	}
	return []any{"error.message", err.Error()}
}

// WithDuration adds a duration in a consistent format (ECS: event.duration)
func WithDuration(d time.Duration) []any {
	return []any{"event.duration", float64(d.Nanoseconds()) / 1e6} // Duration in milliseconds
}

// WithFields adds multiple fields at once
func WithFields(fields map[string]any) []any {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return args
}

// ErrorWithStack logs an error with stack trace (in debug mode)
func ErrorWithStack(msg string, err error, args ...any) {
	if Logger == nil {
		return
	}

	allArgs := enrichArgs(args...)
	allArgs = append(allArgs, "error.message", err.Error())

	// Add stack trace in debug mode
	ctx := context.Background()
	if Logger.Enabled(ctx, slog.LevelDebug) {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		allArgs = append(allArgs, "stack", string(buf[:n]))
	}

	logInternal(ctx, slog.LevelError, msg, allArgs...)
}

// Shutdown gracefully shuts down the logger (useful for async logging)
func Shutdown() {
	if useAsync.Load() && asyncLog != nil {
		asyncLog.Shutdown()
	}
}

// Context key types for request correlation
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
	traceIDKey   contextKey = "trace_id"
)

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext extracts request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extracts user ID from context
func UserIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}

// WithTraceID adds trace ID to context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// TraceIDFromContext extracts trace ID from context
func TraceIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return ""
}
