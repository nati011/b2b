package db

import "errors"

var (
	ErrSysUnknown           = errors.New("unknown error")
	ErrSysNoRows            = errors.New("no rows")
	ErrSysResultSetMismatch = errors.New("result set mismatch")
)
