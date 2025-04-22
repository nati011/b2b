package migration

import (
	"fmt"
	"os"
)

func Core_db_functions() ([]byte, error) {
<<<<<<< HEAD
	filepath := "/home/ruth/Documents/work/nonkifiya/b2b_proj/b2b/migration/core_db_functions.sql"
=======
	filepath := "/home/natanel/personal/b2b_clean/b2b/migration/core_db_functions.sql"
>>>>>>> f7d0812b (+ passing integration tests order)
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}

func Core_db_schema() ([]byte, error) {
<<<<<<< HEAD
	filepath := "/home/ruth/Documents/work/nonkifiya/b2b_proj/b2b/migration/core_db.sql"
=======
	filepath := "/home/natanel/personal/b2b_clean/b2b/migration/core_db.sql"
>>>>>>> f7d0812b (+ passing integration tests order)
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}
