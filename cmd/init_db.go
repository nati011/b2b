package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	ErrFailedToExecuteMigration  = errors.New("oopsy, failed to execute migration")
	ErrFailedToReadMigrationFile = errors.New("oopsy, failed to read migration file")
	ErrFailedToOpenDB            = errors.New("oopsy, failed to open database connection")
	ErrFailedToConnectDB         = errors.New("oopsy, failed to connect with db")
)

func InitDB(connectionString string, file_location string) *sql.DB {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Panic(err.Error())
	}
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Panic(err.Error())
	}

	// ddl
	// err = runMigration(db, file_location+`/core_db.sql`)

	// if err != nil {
	// 	log.Fatalf("Error running migration: %v", err)
	// }

	// // functions
	// err = runMigration(db, file_location+"/core_db_functions.sql")
	// if err != nil {
	// 	log.Fatalf("Error running migration: %v", err)
	// }
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
		log.Panic(err.Error(), filename)
	}
	return nil
}
