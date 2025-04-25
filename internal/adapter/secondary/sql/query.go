package handler

import (
	"context"
	"database/sql"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type QueryResult struct {
	Row  *sql.Row
	Rows *sql.Rows
}

func MustQueryRow(db *sql.DB, ctx context.Context, query string, multiple bool, args ...any) (*QueryResult, error) {
	if multiple {
		rows, err := db.QueryContext(ctx, query, args...)
		if rows.Scan() != nil {
			switch rows.Scan() {
			case sql.ErrNoRows:
				return nil, port_commons.ErrNoRows
			}
		}
		if err != nil {
			return nil, port_commons.ErrSysUnknown
		}
		return &QueryResult{Rows: rows}, nil
	}

	row := db.QueryRowContext(ctx, query, args...)
	if row.Scan() != nil {
		switch row.Scan() {
		case sql.ErrNoRows:
			return nil, port_commons.ErrNoRows
		default:
			return nil, port_commons.ErrSysUnknown
		}
	}

	return &QueryResult{Row: row}, nil
}
