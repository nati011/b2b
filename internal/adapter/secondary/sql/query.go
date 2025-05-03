package handler

import (
	"context"
	"database/sql"
	"log"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type QueryMaster struct {
	args                  []any
	multiRowResultSetDest []any
	singleRowResultSet    []any
	db                    *sql.DB
	ctx                   context.Context
	query                 string
}

type Option func(*QueryMaster)

func NewQuery(options ...Option) QueryMaster {
	svr := QueryMaster{}
	for _, o := range options {
		o(&svr)
	}
	return svr
}

func WithMultiRowResultSet(args []any, dest []any) Option {
	return func(q *QueryMaster) {
		q.args = args
		q.multiRowResultSetDest = dest
	}
}

func WithSingleRowResultSet(args []any, result []any) Option {
	return func(q *QueryMaster) {
		q.args = args
		q.singleRowResultSet = result
	}
}

func WithDB(db *sql.DB) Option {
	return func(q *QueryMaster) {
		q.db = db
	}
}

func WithCtx(ctx context.Context) Option {
	return func(q *QueryMaster) {
		q.ctx = ctx
	}
}

func WithQuery(query string) Option {
	return func(q *QueryMaster) {
		q.query = query
	}
}

// todo: manage offset and limit
func (s QueryMaster) DoMultiQuery() ([][]any, error) {
	rows, err := s.db.QueryContext(s.ctx, s.query, s.args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, port_commons.ErrSysUnknown
	}
	defer rows.Close()

	result := make([][]any, 0)

	for rows.Next() {
		// Create a new slice for each row
		rowValues := make([]any, len(s.multiRowResultSetDest))

		// Initialize each element with a pointer to a new zero value
		for i := range rowValues {
			rowValues[i] = new(any)
		}

		if err := rows.Scan(rowValues...); err != nil {
			log.Printf("unable to scan row: %q", err)
			return nil, port_commons.ErrSysUnknown
		}

		// Dereference the pointers to get the actual values
		finalRow := make([]any, len(rowValues))
		for i, val := range rowValues {
			finalRow[i] = *(val.(*any))
		}

		result = append(result, finalRow)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return nil, port_commons.ErrSysUnknown
	}

	if len(result) == 0 {
		return nil, port_commons.ErrSysNoRows
	}

	return result, nil
}

func (s QueryMaster) DoSingleQuery() error {
	row := s.db.QueryRowContext(s.ctx, s.query, s.args...)
	if row.Err() != nil {
		return port_commons.ErrSysUnknown
	}
	if s.singleRowResultSet != nil {
		err := row.Scan(s.singleRowResultSet...)
		if err != nil {
			return port_commons.ErrSysNoRows
		}
	}
	return nil
}

type QueryResult struct {
	Row  *sql.Row
	Rows *sql.Rows
}

func MustQueryRow(db *sql.DB, ctx context.Context, query string, multiple bool, args ...any) (*QueryResult, error) {
	return nil, nil
}
