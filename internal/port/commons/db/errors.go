package db

import "errors"

var (
	ErrSysUnknown = errors.New("unknown error")
	ErrNoRows     = errors.New("no rows")
)
