package idempotency

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// Record represents a persisted idempotency entry.
type Record struct {
	Key             string
	Method          string
	Path            string
	RequestHash     string
	Processing      bool
	ResponseStatus  sql.NullInt64
	ResponseHeaders map[string][]string
	ResponseBody    []byte
}

// Store persists idempotency executions.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// ClaimState describes the outcome for a key acquisition attempt.
type ClaimState int

const (
	Claimed ClaimState = iota
	Completed
	InProgress
	Conflicted
)

// Claim attempts to register or acquire processing rights for a key.
func (s *Store) Claim(ctx context.Context, key, method, path, requestHash string) (Record, ClaimState, error) {
	var rec Record

	insertQuery := `
		INSERT INTO idempotency_keys (idempotency_key, method, path, request_hash, processing)
		VALUES ($1, $2, $3, $4, TRUE)
		ON CONFLICT (idempotency_key) DO NOTHING
		RETURNING idempotency_key, method, path, request_hash, processing, response_status, response_headers, response_body
	`

	var headersBytes []byte
	err := s.db.QueryRowContext(ctx, insertQuery, key, method, path, requestHash).Scan(
		&rec.Key,
		&rec.Method,
		&rec.Path,
		&rec.RequestHash,
		&rec.Processing,
		&rec.ResponseStatus,
		&headersBytes,
		&rec.ResponseBody,
	)

	switch {
	case err == nil:
		// Newly inserted record; ready to process.
		rec.ResponseHeaders, _ = decodeHeaders(headersBytes)
		return rec, Claimed, nil
	case errors.Is(err, sql.ErrNoRows):
		// Need to load existing record.
	default:
		return Record{}, Claimed, err
	}

	record, err := s.loadRecord(ctx, key)
	if err != nil {
		return Record{}, Claimed, err
	}

	if record.Method != method || record.Path != path || record.RequestHash != requestHash {
		return record, Conflicted, nil
	}

	if record.ResponseStatus.Valid {
		return record, Completed, nil
	}

	if record.Processing {
		return record, InProgress, nil
	}

	// Attempt to acquire processing rights.
	updateQuery := `
		UPDATE idempotency_keys
		SET processing = TRUE, updated_at = $2
		WHERE idempotency_key = $1 AND processing = FALSE AND response_status IS NULL
		RETURNING idempotency_key, method, path, request_hash, processing, response_status, response_headers, response_body
	`

	now := time.Now()
	err = s.db.QueryRowContext(ctx, updateQuery, key, now).Scan(
		&record.Key,
		&record.Method,
		&record.Path,
		&record.RequestHash,
		&record.Processing,
		&record.ResponseStatus,
		&headersBytes,
		&record.ResponseBody,
	)
	if err == nil {
		record.ResponseHeaders, _ = decodeHeaders(headersBytes)
		return record, Claimed, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return record, InProgress, nil
	}
	return Record{}, Claimed, err
}

// SaveResponse persists the HTTP response and releases the key for reuse.
func (s *Store) SaveResponse(ctx context.Context, key string, status int, headers map[string][]string, body []byte) error {
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	query := `
		UPDATE idempotency_keys
		SET processing = FALSE,
		    response_status = $2,
		    response_headers = $3,
		    response_body = $4,
		    updated_at = $5
		WHERE idempotency_key = $1
	`

	_, err = s.db.ExecContext(ctx, query, key, status, headersJSON, body, time.Now())
	return err
}

func (s *Store) loadRecord(ctx context.Context, key string) (Record, error) {
	var rec Record
	var headersBytes []byte

	query := `
		SELECT idempotency_key, method, path, request_hash, processing, response_status, response_headers, response_body
		FROM idempotency_keys
		WHERE idempotency_key = $1
	`

	err := s.db.QueryRowContext(ctx, query, key).Scan(
		&rec.Key,
		&rec.Method,
		&rec.Path,
		&rec.RequestHash,
		&rec.Processing,
		&rec.ResponseStatus,
		&headersBytes,
		&rec.ResponseBody,
	)
	if err != nil {
		return Record{}, err
	}

	rec.ResponseHeaders, _ = decodeHeaders(headersBytes)
	return rec, nil
}

func decodeHeaders(data []byte) (map[string][]string, error) {
	if len(data) == 0 {
		return map[string][]string{}, nil
	}
	var headers map[string][]string
	if err := json.Unmarshal(data, &headers); err != nil {
		return map[string][]string{}, err
	}
	return headers, nil
}
