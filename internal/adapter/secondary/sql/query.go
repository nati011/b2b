package handler

import (
	"context"
	"database/sql"
	"log"

	"b2b.nati011.github.com/config"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type QueryMaster struct {
	hasMultipleResultSet bool
	args                 []any
	multiRowResultSet    [][]any
	singleRowResultSet   []any
	db                   *sql.DB
	ctx                  context.Context
	query                string
	pagination           *config.Pagination
}

type Option func(*QueryMaster)

func NewQuery(options ...Option) QueryMaster {
	svr := QueryMaster{}
	for _, o := range options {
		o(&svr)
	}
	return svr
}

func WithMultiRowResultSet(args []any, results [][]any) Option {
	return func(q *QueryMaster) {
		q.hasMultipleResultSet = true
		q.args = args
		q.multiRowResultSet = results
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

// todo: offset and limit
func (s QueryMaster) DoStuff() error {
	if s.hasMultipleResultSet {
		dest := s.multiRowResultSet[0]
		rows, _ := s.db.QueryContext(s.ctx, s.query, s.args...)
		err := rows.Scan(dest...)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return port_commons.ErrSysNoRows
			}
		}

		defer rows.Close()
		//build list of lists from the one sample on top

		for rows.Next() {
			if err := rows.Scan(dest...); err != nil {
				log.Printf("unable to scan row: %q", err)
				return port_commons.ErrSysUnknown
			}
			s.multiRowResultSet = append(s.multiRowResultSet, dest)
		}

		if err := rows.Err(); err != nil {
			log.Printf("error occurred during rows iteration: %q", err)
			return port_commons.ErrSysUnknown
		}

		return nil
	}

	//single row result set
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
