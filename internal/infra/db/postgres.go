package db

import (
	"context"
	"marketplace/pkg/logger"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewPostgres(conn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", conn)
	if err != nil {
		logger.Error("DB open failed", "error", err)
		return nil, err
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		logger.Error("DB ping failed", "error", err)
		return nil, err
	}

	logger.Info("Database connection established successfully")
	return db, nil
}
