package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ecsServiceName        = "service.name"
	ecsServiceVersion     = "service.version"
	ecsServiceEnvironment = "service.environment"
)

var (
	// cachedHostname caches the hostname to avoid system calls on every log
	cachedHostname string
	hostnameOnce   sync.Once
)

// ecsHandler wraps a slog.JSONHandler and transforms fields to ECS format
type ecsHandler struct {
	handler slog.Handler
	writer  io.Writer
}

// newECSHandler creates a new ECS-compliant handler
func newECSHandler(w io.Writer, opts *slog.HandlerOptions) *ecsHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	// Create a custom ReplaceAttr function that transforms fields to ECS
	ecsOpts := *opts
	originalReplaceAttr := opts.ReplaceAttr
	ecsOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		// Apply original ReplaceAttr if it exists
		if originalReplaceAttr != nil {
			a = originalReplaceAttr(groups, a)
		}

		// Transform standard slog fields to ECS
		if transformed := transformStandardSlogField(groups, a); transformed != nil {
			return *transformed
		}

		// Transform custom fields based on their names
		transformedKey := transformFieldName(a.Key, groups)
		if transformedKey != a.Key {
			return slog.Attr{Key: transformedKey, Value: a.Value}
		}

		return a
	}

	baseHandler := slog.NewJSONHandler(w, &ecsOpts)
	return &ecsHandler{
		handler: baseHandler,
		writer:  w,
	}
}

// transformStandardSlogField transforms standard slog fields to ECS format
// Returns nil if the field is not a standard slog field that needs transformation
func transformStandardSlogField(groups []string, a slog.Attr) *slog.Attr {
	if len(groups) != 0 {
		return nil
	}

	switch a.Key {
	case slog.TimeKey:
		return transformTimeField(a)
	case slog.LevelKey:
		return transformLevelField(a)
	case slog.MessageKey:
		return transformMessageField(a)
	case slog.SourceKey:
		return transformSourceField(a)
	default:
		return nil
	}
}

// transformTimeField transforms time field to @timestamp
func transformTimeField(a slog.Attr) *slog.Attr {
	t, ok := a.Value.Any().(time.Time)
	if !ok {
		return nil
	}
	transformed := slog.String("@timestamp", t.Format(time.RFC3339Nano))
	return &transformed
}

// transformLevelField transforms level field to log.level
func transformLevelField(a slog.Attr) *slog.Attr {
	level := a.Value.String()
	ecsLevel := mapSlogLevelToECS(level)
	transformed := slog.String("log.level", ecsLevel)
	return &transformed
}

// transformMessageField transforms message field
func transformMessageField(a slog.Attr) *slog.Attr {
	return &slog.Attr{Key: "message", Value: a.Value}
}

// transformSourceField transforms source field to log.origin
func transformSourceField(a slog.Attr) *slog.Attr {
	source, ok := a.Value.Any().(slog.Source)
	if !ok {
		return nil
	}
	transformed := slog.Group("log.origin",
		slog.String("file.name", source.File),
		slog.Int("file.line", source.Line),
	)
	return &transformed
}

// mapSlogLevelToECS converts slog level strings to ECS log levels
func mapSlogLevelToECS(level string) string {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return "debug"
	case "info":
		return "info"
	case "warn", "warning":
		return "warn"
	case "error":
		return "error"
	default:
		return "info"
	}
}

// transformFieldName transforms custom field names to ECS field names
func transformFieldName(key string, groups []string) string {
	// If already in a group, preserve the group structure
	if len(groups) > 0 {
		groupPath := strings.Join(groups, ".")
		transformedKey := transformFieldName(key, nil)
		return groupPath + "." + transformedKey
	}

	// Map common field names to ECS equivalents
	fieldMap := map[string]string{
		// Application/Service fields
		"app_name":    ecsServiceName,
		"app_version": ecsServiceVersion,
		"env":         ecsServiceEnvironment,

		// HTTP fields
		"method":        "http.request.method",
		"path":          "url.path",
		"status":        "http.response.status_code",
		"query":         "url.query",
		"response":      "http.response.body.content",
		"response_size": "http.response.body.bytes",

		// User fields
		"user_id": "user.id",

		// Network fields
		"client_ip": "source.ip",

		// Duration
		"duration_ms": "event.duration",
		"duration":    "event.duration", // Keep as string format

		// Error fields
		"error": "error.message",
		"stack": "log.origin.stacktrace",

		// Trace/Request correlation
		"request_id": "trace.id",
		"trace_id":   "trace.id",
	}

	if ecsKey, ok := fieldMap[key]; ok {
		return ecsKey
	}

	// If no mapping found, return as-is (custom fields are allowed in ECS)
	return key
}

// Enabled reports whether the handler handles records at the given level
func (h *ecsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle handles the record
func (h *ecsHandler) Handle(ctx context.Context, r slog.Record) error {
	// The transformation is already handled by ReplaceAttr in the handler options
	// Just pass through to the underlying handler
	return h.handler.Handle(ctx, r)
}

// WithAttrs returns a new handler with the given attributes
func (h *ecsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	transformed := make([]slog.Attr, len(attrs))
	for i, attr := range attrs {
		transformed[i] = transformAttrToECS(attr)
	}
	return &ecsHandler{
		handler: h.handler.WithAttrs(transformed),
		writer:  h.writer,
	}
}

// WithGroup returns a new handler with the given group
func (h *ecsHandler) WithGroup(name string) slog.Handler {
	return &ecsHandler{
		handler: h.handler.WithGroup(name),
		writer:  h.writer,
	}
}

// transformAttrToECS transforms a single attribute to ECS format
func transformAttrToECS(a slog.Attr) slog.Attr {
	// Handle nested groups
	if a.Value.Kind() == slog.KindGroup {
		groupAttrs := a.Value.Group()
		transformedKey := transformFieldName(a.Key, nil)
		// Convert []slog.Attr to []any for slog.Group
		groupArgs := make([]any, 0, len(groupAttrs)*2)
		for _, ga := range groupAttrs {
			transformedGA := transformAttrToECS(ga)
			groupArgs = append(groupArgs, transformedGA.Key, transformedGA.Value.Any())
		}
		return slog.Group(transformedKey, groupArgs...)
	}

	// Transform the key
	transformedKey := transformFieldName(a.Key, nil)
	if transformedKey != a.Key {
		return slog.Attr{Key: transformedKey, Value: a.Value}
	}

	return a
}

// addECSFields adds standard ECS fields that should be present in all logs
func addECSFields(attrs []any, appCfg map[string]any) []any {
	// Add service fields if available (check ECS names first, then fall back to old names)
	var serviceName, serviceVersion, serviceEnv string

	if name, ok := appCfg[ecsServiceName].(string); ok {
		serviceName = name
	} else if name, ok := appCfg["app_name"].(string); ok {
		serviceName = name
	}
	if serviceName != "" {
		attrs = append(attrs, ecsServiceName, serviceName)
	}

	if version, ok := appCfg[ecsServiceVersion].(string); ok {
		serviceVersion = version
	} else if version, ok := appCfg["app_version"].(string); ok {
		serviceVersion = version
	}
	if serviceVersion != "" {
		attrs = append(attrs, ecsServiceVersion, serviceVersion)
	}

	if env, ok := appCfg[ecsServiceEnvironment].(string); ok {
		serviceEnv = env
	} else if env, ok := appCfg["env"].(string); ok {
		serviceEnv = env
	}
	if serviceEnv != "" {
		attrs = append(attrs, ecsServiceEnvironment, serviceEnv)
	}

	// Add hostname if available (cached to avoid system calls)
	hostnameOnce.Do(func() {
		if hostname, err := os.Hostname(); err == nil {
			cachedHostname = hostname
		}
	})
	if cachedHostname != "" {
		attrs = append(attrs, "host.name", cachedHostname)
	}

	// Add ECS version
	attrs = append(attrs, "ecs.version", "8.11.0")

	return attrs
}
