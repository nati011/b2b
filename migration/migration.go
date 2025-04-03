package migration

import (
	"fmt"
	"os"
)

func Core_db_functions() ([]byte, error) {
	filepath := "../core_db_functions.sql"
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}

func Core_db_schema() ([]byte, error) {
	filepath := "../core_db.sql"
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}
