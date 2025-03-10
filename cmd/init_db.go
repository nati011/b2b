package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	ErrFailedToExecuteMigration  = errors.New("oopsy, failed to execute migration")
	ErrFailedToReadMigrationFile = errors.New("oopsy, failed to read migration file")
	ErrFailedToOpenDB            = errors.New("oopsy, failed to open database connection")
	ErrFailedToConnectDB         = errors.New("oopsy, failed to connect with db")
)

func InitDB(connectionString string) *sql.DB {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Panic(ErrFailedToOpenDB)
	}
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Panic(ErrFailedToConnectDB)
	}

	// ddl
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db.sql")
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}

	// functions
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db_functions.sql")
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}
	return db
}

func runMigration(db *sql.DB, filename string) error {
	// Read
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Panic(ErrFailedToReadMigrationFile)
	}

	// Execute
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Panic(ErrFailedToExecuteMigration)
	}

	return nil
}
