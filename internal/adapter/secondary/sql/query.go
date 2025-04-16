package handler

import (
	"context"
	"database/sql"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

func MustQueryRow(db *sql.DB, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, port_commons.ErrNoRows
		default:
			return nil, port_commons.ErrSysUnknown
		}
	}
	return rows, nil
}
