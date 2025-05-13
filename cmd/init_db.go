package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"

	"b2b.nati011.github.com/config"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	ErrFailedToExecuteMigration  = errors.New("oopsy, failed to execute migration")
	ErrFailedToReadMigrationFile = errors.New("oopsy, failed to read migration file")
	ErrFailedToOpenDB            = errors.New("oopsy, failed to open database connection")
	ErrFailedToConnectDB         = errors.New("oopsy, failed to connect with db")
)

func InitDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("pgx", cfg.CoreDBConnectionString)
	if err != nil {
		log.Panic(err.Error())
	}
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Panic(err.Error())
	}

	switch cfg.Env {
	case "development":
		InitDBDevelopment(db, cfg.FileLocation, cfg)
	case "staging":
		InitDBStaging(db, cfg.FileLocation, cfg)
	case "production":
		InitDBProduction(db, cfg.FileLocation, cfg)
	}

	return db
}

func InitDBDevelopment(db *sql.DB, file_location string, cfg *config.Config) {
	// ddl
	err := runMigration(db, file_location+`/core_db.sql`)

	if err != nil {
		log.Fatalf("Error running ddl migration: %v", err)
	}

	// functions
	err = runMigration(db, file_location+"/core_db_functions.sql")
	if err != nil {
		log.Fatalf("Error running function migration: %v", err)
	}

	// seed
	err = runMigration(db, file_location+"/core_init_migration_script.sql")
	if err != nil {
		log.Fatalf("Error running seed migration: %v", err)
	}
}

func InitDBStaging(db *sql.DB, file_location string, cfg *config.Config) {
}

func InitDBProduction(db *sql.DB, file_location string, cfg *config.Config) {
	// // ddl
	// err := runMigration(db, file_location+`/core_db.sql`)

	// if err != nil {
	// 	log.Fatalf("Error running ddl migration: %v", err)
	// }

	// // functions
	// err = runMigration(db, file_location+"/core_db_functions.sql")
	// if err != nil {
	// 	log.Fatalf("Error running function migration: %v", err)
	// }

	// // seed
	// err = runMigration(db, file_location+"/core_init_migration_script.sql")
	// if err != nil {
	// 	log.Fatalf("Error running seed migration: %v", err)
	// }
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
