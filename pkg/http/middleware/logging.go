package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/logger"
)

const (
	maxResponseBodyLength = 500
)

type LoggingMiddleware struct {
}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

type statusRecorder struct {
	http.ResponseWriter
	status      int
	response    *bytes.Buffer
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(status int) {
	if !r.wroteHeader {
		r.status = status
		r.wroteHeader = true
		r.ResponseWriter.WriteHeader(status)
	}
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	// Always capture the response body, even if empty (to track writes)
	r.response.Write(b)
	n, err := r.ResponseWriter.Write(b)
	return n, err
}

// Flush implements http.Flusher to ensure buffered data is captured
func (r *statusRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (lm *LoggingMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseBuffer := &bytes.Buffer{}
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
			response:       responseBuffer,
			wroteHeader:    false,
		}

		defer func() {
			if logInfoEnabled() {
				lm.logRequest(r, recorder, start)
			}
		}()

		next.ServeHTTP(recorder, r)
	})
}

// logRequest logs the HTTP request with ECS-compliant fields
func (lm *LoggingMiddleware) logRequest(r *http.Request, recorder *statusRecorder, start time.Time) {
	duration := time.Since(start)
	clientIP := getClientIP(r)
	queryString := r.URL.RawQuery
	fullBody := recorder.response.String()
	responseBody := truncateString(fullBody, maxResponseBodyLength)

	msg := buildLogMessage(r, queryString, recorder.status)
	attrs := buildBaseAttributes(r, recorder, fullBody, duration, clientIP)
	attrs = addOptionalAttributes(attrs, r, queryString)
	attrs = addResponseBodyAttribute(attrs, fullBody, responseBody, recorder.status)

	logger.InfoContext(r.Context(), msg, attrs...)
}

// buildLogMessage creates a readable log message from the request
func buildLogMessage(r *http.Request, queryString string, status int) string {
	msg := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	if queryString != "" {
		msg += "?" + queryString
	}
	msg += fmt.Sprintf(" -> %d", status)
	return msg
}

// buildBaseAttributes creates the base ECS-compliant attributes
func buildBaseAttributes(r *http.Request, recorder *statusRecorder, fullBody string, duration time.Duration, clientIP string) []any {
	return []any{
		// HTTP request fields
		"http.request.method", r.Method,
		"url.path", r.URL.Path,
		"url.full", r.URL.String(),

		// HTTP response fields
		"http.response.status_code", recorder.status,
		"http.response.body.bytes", len(fullBody),

		// Event fields
		"event.duration", duration.Nanoseconds() / 1e6, // milliseconds
		"event.category", "web",
		"event.action", "http_request",
		"event.kind", "event",
		"event.type", []string{"access"},

		// Network fields
		"source.ip", clientIP,

		// HTTP version
		"http.version", r.Proto,
	}
}

// addOptionalAttributes adds optional attributes like user ID, query string, and user agent
func addOptionalAttributes(attrs []any, r *http.Request, queryString string) []any {
	// Extract user ID from context if available (ECS: user.id)
	if user := httputil.UserFromContext(r.Context()); user != nil {
		attrs = append(attrs, "user.id", user.ID)
	}

	// URL query string
	if queryString != "" {
		attrs = append(attrs, "url.query", queryString)
	}

	// Add user agent
	if ua := r.Header.Get("User-Agent"); ua != "" {
		attrs = append(attrs, "user_agent.original", ua)
	}

	return attrs
}

// addResponseBodyAttribute determines and adds the response body content attribute
func addResponseBodyAttribute(attrs []any, fullBody, responseBody string, status int) []any {
	responseValue := determineResponseValue(fullBody, responseBody, status)
	return append(attrs, "http.response.body.content", responseValue)
}

// determineResponseValue determines the appropriate response value to log
func determineResponseValue(fullBody, responseBody string, status int) string {
	if fullBody == "" {
		return "(empty)"
	}

	if summary := extractResponseSummary(fullBody); summary != "" {
		return summary
	}

	if status >= 400 || len(fullBody) <= maxResponseBodyLength {
		// For errors or small responses, show the full body
		return fullBody
	}

	// For large successful responses, show preview
	return responseBody + " [truncated]"
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// extractResponseSummary extracts key information from JSON responses for readable logging
func extractResponseSummary(body string) string {
	if body == "" {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// Not JSON, return as-is if small
		if len(body) <= 100 {
			return body
		}
		return ""
	}

	parts := extractCommonFields(data)
	parts = append(parts, extractItemFields(data)...)

	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}

	return extractCompactJSON(data)
}

// extractCommonFields extracts common pagination and error fields
func extractCommonFields(data map[string]interface{}) []string {
	var parts []string

	if items, ok := data["items"].([]interface{}); ok {
		parts = append(parts, fmt.Sprintf("items=%d", len(items)))
	}
	if total, ok := data["total"].(float64); ok {
		parts = append(parts, fmt.Sprintf("total=%.0f", total))
	}
	if page, ok := data["page"].(float64); ok {
		parts = append(parts, fmt.Sprintf("page=%.0f", page))
	}
	if limit, ok := data["limit"].(float64); ok {
		parts = append(parts, fmt.Sprintf("limit=%.0f", limit))
	}
	if err, ok := data["error"].(string); ok {
		parts = append(parts, fmt.Sprintf("error=%s", err))
	}
	if msg, ok := data["message"].(string); ok {
		parts = append(parts, fmt.Sprintf("message=%s", msg))
	}
	if id, ok := data["id"].(string); ok && len(id) > 0 {
		parts = append(parts, fmt.Sprintf("id=%s", id))
	}

	return parts
}

// extractItemFields extracts fields commonly found in single item responses
func extractItemFields(data map[string]interface{}) []string {
	var parts []string

	if email, ok := data["email"].(string); ok {
		parts = append(parts, fmt.Sprintf("email=%s", email))
	}
	if name, ok := data["name"].(string); ok {
		parts = append(parts, fmt.Sprintf("name=%s", name))
	}
	if status, ok := data["status"].(string); ok {
		parts = append(parts, fmt.Sprintf("status=%s", status))
	}

	return parts
}

// extractCompactJSON returns compact JSON as fallback if within size limit
func extractCompactJSON(data map[string]interface{}) string {
	compact, err := json.Marshal(data)
	if err == nil && len(compact) <= maxResponseBodyLength {
		return string(compact)
	}
	return ""
}

func logInfoEnabled() bool {
	if logger.Logger == nil {
		return false
	}
	// Use FromContext as an example of the alternative logger function
	ctxLogger := logger.FromContext(context.Background())
	if ctxLogger == nil {
		return false
	}
	return ctxLogger.Enabled(context.Background(), slog.LevelInfo)
}
