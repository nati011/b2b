package migration

import (
	"fmt"
	"os"
)

func Core_db_functions() ([]byte, error) {
<<<<<<< HEAD
	filepath := "/home/natanel/personal/b2b_clean/b2b/migration/core_db_functions.sql"
=======
	filepath := "/home/ruth/Documents/work/nonkifiya/b2b_proj/b2b/migration/core_db_functions.sql"
>>>>>>> c481e966 (init handle multiple payment gateway)
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}

func Core_db_schema() ([]byte, error) {
<<<<<<< HEAD
	filepath := "/home/natanel/personal/b2b_clean/b2b/migration/core_db.sql"
=======
	filepath := "/home/ruth/Documents/work/nonkifiya/b2b_proj/b2b/migration/core_db.sql"
>>>>>>> c481e966 (init handle multiple payment gateway)
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return sqlBytes, nil
}
