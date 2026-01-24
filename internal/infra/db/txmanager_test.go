package db

import (
	"context"
	"marketplace/test"
	"database/sql"
	"errors"
	"testing"
)

func newTestTxManager(t *testing.T) (*sql.DB, TxManager, func()) {
	t.Helper()

	dbConn, cleanup := test.SetupTestDB(t)
	txMgr := NewTxManager(dbConn)

	return dbConn, txMgr, cleanup
}

func TestTxManagerCommitOnSuccess(t *testing.T) {
	dbConn, txMgr, cleanup := newTestTxManager(t)
	defer cleanup()

	ctx := context.Background()
	if _, err := dbConn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS tx_test (id SERIAL PRIMARY KEY, value INT NOT NULL)`); err != nil {
		t.Fatalf("failed to create tx_test table: %v", err)
	}

	err := txMgr.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO tx_test (value) VALUES (1)`); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithinTx returned error: %v", err)
	}

	var count int
	if err := dbConn.QueryRowContext(ctx, `SELECT COUNT(*) FROM tx_test WHERE value = 1`).Scan(&count); err != nil {
		t.Fatalf("failed to query tx_test: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row after commit, got %d", count)
	}
}

func TestTxManagerRollbackOnError(t *testing.T) {
	dbConn, txMgr, cleanup := newTestTxManager(t)
	defer cleanup()

	ctx := context.Background()
	if _, err := dbConn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS tx_test_err (id SERIAL PRIMARY KEY, value INT NOT NULL)`); err != nil {
		t.Fatalf("failed to create tx_test_err table: %v", err)
	}

	testErr := errors.New("forced error")
	err := txMgr.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO tx_test_err (value) VALUES (1)`); err != nil {
			return err
		}
		return testErr
	})
	if !errors.Is(err, testErr) {
		t.Fatalf("expected error %v, got %v", testErr, err)
	}

	var count int
	if err := dbConn.QueryRowContext(ctx, `SELECT COUNT(*) FROM tx_test_err WHERE value = 1`).Scan(&count); err != nil {
		t.Fatalf("failed to query tx_test_err: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 rows after rollback, got %d", count)
	}
}

func TestTxManagerPanicRollsBack(t *testing.T) {
	dbConn, txMgr, cleanup := newTestTxManager(t)
	defer cleanup()

	ctx := context.Background()
	if _, err := dbConn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS tx_test_panic (id SERIAL PRIMARY KEY, value INT NOT NULL)`); err != nil {
		t.Fatalf("failed to create tx_test_panic table: %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic to be re-thrown from WithinTx")
		}

		var count int
		if err := dbConn.QueryRowContext(ctx, `SELECT COUNT(*) FROM tx_test_panic WHERE value = 1`).Scan(&count); err != nil {
			t.Fatalf("failed to query tx_test_panic: %v", err)
		}
		if count != 0 {
			t.Fatalf("expected 0 rows after panic rollback, got %d", count)
		}
	}()

	_ = txMgr.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO tx_test_panic (value) VALUES (1)`); err != nil {
			return err
		}
		panic("boom")
	})
}
