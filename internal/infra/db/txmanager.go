package db

import (
	"context"
	"database/sql"
)

// TxManager defines a helper for running code within a database transaction.
//
// It is infrastructure-level and should be used by repositories or services that
// require atomic multi-statement operations.
type TxManager interface {
	// WithinTx begins a transaction, executes fn, and:
	//   - commits if fn returns nil,
	//   - rolls back if fn returns an error,
	//   - rolls back and re-panics if fn panics.
	WithinTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}

// txManager is the default TxManager implementation backed by *sql.DB.
type txManager struct {
	db *sql.DB
}

// NewTxManager constructs a new TxManager for the provided database connection.
func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db: db}
}

// WithinTx executes fn within a database transaction and ensures proper commit /
// rollback semantics, including panic safety.
func (m *txManager) WithinTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) (err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}

		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	err = fn(ctx, tx)
	return
}
