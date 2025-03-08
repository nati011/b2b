package main

import (
	"context"
	"database/sql"
	"log"
)

func InitDB(connectionString string) *sql.DB {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Panic("failed to open database connection")
	}
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Panic("failed to connect to database")
	}
	return db
}
