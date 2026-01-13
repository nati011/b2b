package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func getMigrationPath(filename string) (string, error) {
	// Get the directory where this source file is located
	_, currentFile, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("could not get current file path")
	}
	
	// Get the migration directory (same directory as this file)
	migrationDir := filepath.Dir(currentFile)
	filePath := filepath.Join(migrationDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}
	
	// Fallback: try relative to working directory
	wd, err := os.Getwd()
	if err == nil {
		possiblePaths := []string{
			filepath.Join(wd, "migration", filename),
			filepath.Join(wd, "..", "migration", filename),
		}
		
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}
	
	return "", fmt.Errorf("could not find migration file %s (searched: %s and relative paths)", filename, filePath)
}

func Core_db_functions() ([]byte, error) {
	filepath, err := getMigrationPath("core_db_functions.sql")
	if err != nil {
		return nil, err
	}
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}

func Core_db_schema() ([]byte, error) {
	filepath, err := getMigrationPath("core_db.sql")
	if err != nil {
		return nil, err
	}
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}
