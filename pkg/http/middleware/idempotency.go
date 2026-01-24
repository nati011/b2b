package middleware

import (
	"bufio"
	"bytes"
	"context"
	"marketplace/internal/infra/idempotency"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

const idempotencyHeader = "Idempotency-Key"

// IdempotencyMiddleware enforces infrastructure-level idempotency.
type IdempotencyMiddleware struct {
	store *idempotency.Store
}

func NewIdempotencyMiddleware(store *idempotency.Store) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		store: store,
	}
}

func (m *IdempotencyMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.shouldApply(r) {
			next.ServeHTTP(w, r)
			return
		}

		key := strings.TrimSpace(r.Header.Get(idempotencyHeader))
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		hash := hashPayload(r.Method, r.URL.Path, bodyBytes)

		record, state, err := m.store.Claim(r.Context(), key, r.Method, r.URL.Path, hash)
		if err != nil {
			http.Error(w, "failed to process idempotency key", http.StatusInternalServerError)
			return
		}

		switch state {
		case idempotency.Conflicted:
			http.Error(w, "idempotency key reuse with different payload", http.StatusConflict)
			return
		case idempotency.Completed:
			m.replayResponse(w, record)
			return
		case idempotency.InProgress:
			http.Error(w, "idempotent request is still processing", http.StatusConflict)
			return
		case idempotency.Claimed:
			// Continue to handler.
		default:
			http.Error(w, "invalid idempotency state", http.StatusInternalServerError)
			return
		}

		rec := newResponseRecorder(w)
		next.ServeHTTP(rec, r)

		if err := m.persistResponse(r.Context(), key, rec); err != nil {
			// best-effort logging via stdlib since logger may not be configured yet
			fmt.Println("idempotency save error:", err)
		}
	})
}

func (m *IdempotencyMiddleware) shouldApply(r *http.Request) bool {
	return r.Method == http.MethodPost || r.Method == http.MethodPut
}

func hashPayload(method, path string, body []byte) string {
	sum := sha256.Sum256(append([]byte(method+"|"+path+"|"), body...))
	return hex.EncodeToString(sum[:])
}

func (m *IdempotencyMiddleware) replayResponse(w http.ResponseWriter, record idempotency.Record) {
	for key, values := range record.ResponseHeaders {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	status := int(record.ResponseStatus.Int64)
	w.WriteHeader(status)
	if len(record.ResponseBody) > 0 {
		_, _ = w.Write(record.ResponseBody)
	}
}

func (m *IdempotencyMiddleware) persistResponse(ctx context.Context, key string, rec *responseRecorder) error {
	if rec.status == 0 {
		rec.status = http.StatusOK
	}
	headers := map[string][]string{}
	for k, values := range rec.Header() {
		headers[k] = append([]string(nil), values...)
	}
	return m.store.SaveResponse(ctx, key, rec.status, headers, rec.body.Bytes())
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Ensure responseRecorder implements http.Hijacker if the underlying writer supports it.
func (r *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijacker not supported")
	}
	return h.Hijack()
}
